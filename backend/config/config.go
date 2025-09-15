// config/config.go
package config

import (
	"fmt"
	"os"
)

// DatabaseConfig содержит настройки для подключения к базе данных
type DatabaseConfig struct {
	Host     string
	Port     interface{} // может быть string или int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Config содержит все настройки приложения
type Config struct {
	Database DatabaseConfig
	Port     string
}

// LoadConfig загружает конфигурацию из переменных окружения
func LoadConfig() *Config {
	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "todo"),
			Password: getEnv("DB_PASSWORD", "todo"),
			DBName:   getEnv("DB_NAME", "todo"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Port: getEnv("APP_PORT", "8080"),
	}
}

// GetDSN возвращает строку подключения к PostgreSQL
func (c *DatabaseConfig) GetDSN() string {
	var port string
	switch p := c.Port.(type) {
	case int:
		port = fmt.Sprintf("%d", p)
	case string:
		port = p
	default:
		port = "5432"
	}

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// getEnv получает значение переменной окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
