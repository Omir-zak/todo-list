// database/database.go
package database

import (
	"fmt"
	"log"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"todo-list/backend/config"
	"todo-list/backend/internal/models"
)

// Database структура для работы с базой данных
type Database struct {
	DB *gorm.DB
}

// NewConnection создает новое подключение к базе данных (для совместимости с start.go)
func NewConnection(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	var port string
	if p, ok := cfg.Port.(int); ok {
		port = strconv.Itoa(p)
	} else if p, ok := cfg.Port.(string); ok {
		port = p
	} else {
		port = "5432"
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Проверяем соединение
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to database")
	return db, nil
}

// NewDatabase создает новое подключение к базе данных
func NewDatabase(cfg *config.Config) (*Database, error) {
	db, err := NewConnection(&cfg.Database)
	if err != nil {
		return nil, err
	}
	return &Database{DB: db}, nil
}

// Migrate выполняет миграции базы данных
func Migrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	err := db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Todo{},
	)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// Close закрывает соединение с базой данных
func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Seed заполняет базу данных начальными данными
func (d *Database) Seed() error {
	log.Println("Seeding database...")

	// Создаем категории по умолчанию
	categories := []models.Category{
		{Name: "Работа", Color: "#dc3545"},
		{Name: "Личное", Color: "#28a745"},
		{Name: "Покупки", Color: "#ffc107"},
		{Name: "Учеба", Color: "#007bff"},
	}

	for _, category := range categories {
		var existingCategory models.Category
		result := d.DB.Where("name = ?", category.Name).First(&existingCategory)
		if result.Error == gorm.ErrRecordNotFound {
			if err := d.DB.Create(&category).Error; err != nil {
				return fmt.Errorf("failed to create category %s: %w", category.Name, err)
			}
			log.Printf("Created category: %s", category.Name)
		}
	}

	// Создаем пользователя по умолчанию
	defaultUser := models.User{
		Username: "admin",
		Email:    "admin@todo.com",
		Password: "admin123", // В реальном приложении нужно хешировать пароль
	}

	var existingUser models.User
	result := d.DB.Where("username = ?", defaultUser.Username).First(&existingUser)
	if result.Error == gorm.ErrRecordNotFound {
		if err := d.DB.Create(&defaultUser).Error; err != nil {
			return fmt.Errorf("failed to create default user: %w", err)
		}
		log.Printf("Created default user: %s", defaultUser.Username)
	}

	log.Println("Database seeding completed successfully")
	return nil
}
