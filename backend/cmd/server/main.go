package main

import (
	"flag"
	"log"
	"os"

	"github.com/heartmarshall/my-english/internal/app"
	"github.com/heartmarshall/my-english/internal/config"
)

func main() {
	// Флаг для указания пути к YAML конфигу
	configPath := flag.String("config", "", "path to YAML config file (optional)")
	flag.Parse()

	// Проверяем переменную окружения CONFIG_PATH, если флаг не указан
	if *configPath == "" {
		if configPathEnv := os.Getenv("CONFIG_PATH"); configPathEnv != "" {
			*configPath = configPathEnv
		}
	}

	// Загружаем конфигурацию
	// Приоритет: переменные окружения > YAML файл > дефолтные значения
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Создаём приложение
	application, err := app.New(*cfg)
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
