package main

import (
	"log"

	"github.com/heartmarshall/my-english/internal/app"
)

func main() {
	// Загружаем конфигурацию
	cfg := app.LoadFromEnv()

	// Создаём приложение
	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}

	// Запускаем
	if err := application.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
