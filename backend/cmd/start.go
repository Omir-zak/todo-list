package cmd

import (
	"context"
	"embed"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"log"
	"os"
	"strconv"
	"todo-list/backend/config"
	"todo-list/backend/database"
	"todo-list/backend/internal/repository"
	"todo-list/backend/internal/service"
	wailsbind "todo-list/backend/internal/wails"
)

// 🇧 ВАЖНО: путь для go:embed должен указывать на фактическое расположение dist.
// Если этот файл лежит в cmd/desktop, а папка frontend — в корне репо,
// то используйте относительный путь с выходом на два уровня вверх:
//
//go:embed ../../frontend/dist/*
var assets embed.FS

func Start() {
	// 1) Конфиг из ENV (замени на свой loader при необходимости)
	dbCfg := config.DatabaseConfig{
		Host:     getenv("DB_HOST", "localhost"),
		Port:     getenvInt("DB_PORT", 5432),
		User:     getenv("DB_USER", "postgres"),
		Password: getenv("DB_PASSWORD", "postgres"),
		DBName:   getenv("DB_NAME", "todo"),
		SSLMode:  getenv("DB_SSLMODE", "disable"),
	}

	// 2) Инициализация зависимостей ДО запуска Wails (чтобы не городить lazy-прокси)
	db, err := database.NewConnection(&dbCfg)
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatalf("migrate error: %v", err)
	}

	taskRepo := repository.NewTaskRepository(db) // твоя реализация
	taskSvc := service.NewTaskService(taskRepo)  // твоя реализация
	taskAPI := wailsbind.NewTaskAPI(taskSvc)

	// 3) Запуск Wails-приложения
	err = wails.Run(&options.App{
		Title:  "To-Do List (Wails)",
		Width:  1024,
		Height: 700,
		AssetServer: &assetserver.Options{
			Assets: assets, // берёт статические файлы из frontend/dist
		},
		OnStartup: func(ctx context.Context) {
			// если надо — можно сделать что-то при старте окна
		},
		OnShutdown: func(ctx context.Context) {
			// ресурсы уже закрываем через defer, но тут тоже можно что-то добить
		},
		Bind: []interface{}{
			taskAPI, // биндим экспортируемый API в JS-слой Wails
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}

// --- utils ---

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func getenvInt(k string, def int) int {
	if v := os.Getenv(k); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}
