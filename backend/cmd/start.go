package cmd

import (
	"context"
	"embed"
	"log"
	"os"
	"runtime/debug"
	"todo-list/backend"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

func Start(assets embed.FS) {
	// в самом начале main() / Start()
	f, _ := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	defer f.Close()
	log.SetOutput(f)
	log.Println("app starting...")

	defer func() {
		if r := recover(); r != nil {
			log.Printf("PANIC: %v\n", r)
			log.Printf("stack:\n%s", debug.Stack())
		}
	}()

	// Создаем простое приложение без базы данных
	app := backend.NewApp()

	// Запуск Wails-приложения
	err := wails.Run(&options.App{
		Title:  "To-Do List (Wails)",
		Width:  1024,
		Height: 700,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: func(ctx context.Context) {
			// Инициализация при необходимости
		},
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
