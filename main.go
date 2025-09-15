package main

import (
	"embed"
	"todo-list/backend/cmd"
)

//go:embed all:frontend/dist
var Assets embed.FS

func main() {
	cmd.Start(Assets)
}
