package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

// LoadFromYAML загружает конфигурацию из YAML файла.
// Если файл не указан или не существует, возвращает конфигурацию по умолчанию.
// Переменные окружения имеют приоритет над значениями из YAML.
func LoadFromYAML(path string) (Config, error) {
	cfg := DefaultConfig()

	if path == "" {
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Файл не существует, используем дефолтные значения
			return cfg, nil
		}
		return cfg, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("failed to parse YAML config: %w", err)
	}

	// Применяем переменные окружения поверх YAML конфига
	applyEnvOverrides(&cfg)

	return cfg, nil
}

// LoadFromEnv загружает конфигурацию из переменных окружения.
// Использует дефолтные значения, если переменная не установлена.
func LoadFromEnv() Config {
	cfg := DefaultConfig()
	applyEnvOverrides(&cfg)
	return cfg
}

// applyEnvOverrides применяет переменные окружения к конфигурации.
func applyEnvOverrides(cfg *Config) {
	// Server
	if host := os.Getenv("SERVER_HOST"); host != "" {
		cfg.Server.Host = host
	}
	if port := os.Getenv("SERVER_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Server.Port = p
		}
	}
	if timeout := os.Getenv("SERVER_REQUEST_TIMEOUT"); timeout != "" {
		if d, err := time.ParseDuration(timeout); err == nil {
			cfg.Server.RequestTimeout = d
		}
	}

	// Database
	if host := os.Getenv("DB_HOST"); host != "" {
		cfg.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Database.Port = p
		}
	}
	if user := os.Getenv("DB_USER"); user != "" {
		cfg.Database.User = user
	}
	if pass := os.Getenv("DB_PASSWORD"); pass != "" {
		cfg.Database.Password = pass
	}
	if db := os.Getenv("DB_NAME"); db != "" {
		cfg.Database.Database = db
	}
	if ssl := os.Getenv("DB_SSLMODE"); ssl != "" {
		cfg.Database.SSLMode = ssl
	}

	// GraphQL
	if env := os.Getenv("GRAPHQL_PLAYGROUND"); env == "false" {
		cfg.GraphQL.EnablePlayground = false
	}
	if env := os.Getenv("GRAPHQL_INTROSPECTION"); env == "false" {
		cfg.GraphQL.EnableIntrospection = false
	}

	// Logging
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		cfg.Log.Level = LogLevel(level)
	}
	if format := os.Getenv("LOG_FORMAT"); format != "" {
		cfg.Log.Format = format
	}
}
