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

	// Загружаем конфигурацию
	var cfg config.Config
	var err error
	if *configPath != "" {
		cfg, err = config.LoadFromYAML(*configPath)
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}
	} else {
		// Проверяем переменную окружения CONFIG_PATH
		if configPathEnv := os.Getenv("CONFIG_PATH"); configPathEnv != "" {
			cfg, err = config.LoadFromYAML(configPathEnv)
			if err != nil {
				log.Fatalf("Failed to load config: %v", err)
			}
		} else {
			// Используем переменные окружения
			cfg = config.LoadFromEnv()
		}
	}

	// Создаём приложение
	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
