// database/database.go
package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"todo-list/backend/config"
	"todo-list/backend/internal/models"
)

// Database структура для работы с базой данных
type Database struct {
	DB *sql.DB
}

// NewConnection создает новое подключение к базе данных
func NewConnection(cfg *config.DatabaseConfig) (*sql.DB, error) {
	var port string
	if p, ok := cfg.Port.(int); ok {
		port = fmt.Sprintf("%d", p)
	} else if p, ok := cfg.Port.(string); ok {
		port = p
	} else {
		port = "5432"
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
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
func Migrate(db *sql.DB) error {
	log.Println("Running database migrations...")

	// Создание таблицы категорий
	categoryTableSQL := `
	CREATE TABLE IF NOT EXISTS categories (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		color VARCHAR(7) DEFAULT '#007bff',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	// Создание таблицы задач
	todoTableSQL := `
	CREATE TABLE IF NOT EXISTS todos (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		completed BOOLEAN DEFAULT FALSE,
		priority VARCHAR(10) DEFAULT 'medium',
		due_date TIMESTAMP,
		category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	// Создание индексов
	indexesSQL := []string{
		`CREATE INDEX IF NOT EXISTS idx_todos_category_id ON todos(category_id)`,
		`CREATE INDEX IF NOT EXISTS idx_todos_completed ON todos(completed)`,
		`CREATE INDEX IF NOT EXISTS idx_todos_due_date ON todos(due_date)`,
	}

	// Выполняем миграции
	tables := []string{categoryTableSQL, todoTableSQL}
	for _, tableSQL := range tables {
		if _, err := db.Exec(tableSQL); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	// Создаем индексы
	for _, indexSQL := range indexesSQL {
		if _, err := db.Exec(indexSQL); err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// Close закрывает соединение с базой данных
func (d *Database) Close() error {
	return d.DB.Close()
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
		var count int
		err := d.DB.QueryRow("SELECT COUNT(*) FROM categories WHERE name = $1", category.Name).Scan(&count)
		if err != nil {
			return fmt.Errorf("failed to check category existence: %w", err)
		}

		if count == 0 {
			_, err = d.DB.Exec(
				"INSERT INTO categories (name, color) VALUES ($1, $2)",
				category.Name, category.Color,
			)
			if err != nil {
				return fmt.Errorf("failed to create category %s: %w", category.Name, err)
			}
			log.Printf("Created category: %s", category.Name)
		}
	}

	log.Println("Database seeding completed successfully")
	return nil
}
