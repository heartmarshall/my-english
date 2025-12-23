package app

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config — конфигурация приложения.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	GraphQL  GraphQLConfig
}

// ServerConfig — конфигурация HTTP сервера.
type ServerConfig struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DatabaseConfig — конфигурация базы данных.
type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Database        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// DSN возвращает строку подключения к PostgreSQL.
func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode,
	)
}

// GraphQLConfig — конфигурация GraphQL.
type GraphQLConfig struct {
	EnablePlayground    bool
	EnableIntrospection bool
	QueryCacheSize      int
}

// DefaultConfig возвращает конфигурацию по умолчанию.
func DefaultConfig() Config {
	return Config{
		Server: ServerConfig{
			Host:         "0.0.0.0",
			Port:         8080,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
		},
		Database: DatabaseConfig{
			Host:            "localhost",
			Port:            5432,
			User:            "postgres",
			Password:        "postgres",
			Database:        "my_english",
			SSLMode:         "disable",
			MaxOpenConns:    25,
			MaxIdleConns:    5,
			ConnMaxLifetime: 5 * time.Minute,
		},
		GraphQL: GraphQLConfig{
			EnablePlayground:    true,
			EnableIntrospection: true,
			QueryCacheSize:      1000,
		},
	}
}

// LoadFromEnv загружает конфигурацию из переменных окружения.
// Использует дефолтные значения, если переменная не установлена.
func LoadFromEnv() Config {
	cfg := DefaultConfig()

	// Server
	if host := os.Getenv("SERVER_HOST"); host != "" {
		cfg.Server.Host = host
	}
	if port := os.Getenv("SERVER_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Server.Port = p
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

	return cfg
}

// Addr возвращает адрес сервера.
func (c ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
