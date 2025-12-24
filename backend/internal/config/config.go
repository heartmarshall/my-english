package config

import (
	"fmt"
	"time"
)

// Config — конфигурация приложения.
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	GraphQL  GraphQLConfig  `yaml:"graphql"`
	Log      LogConfig      `yaml:"log"`
}

// ServerConfig — конфигурация HTTP сервера.
type ServerConfig struct {
	Host           string        `yaml:"host"`
	Port           int           `yaml:"port"`
	ReadTimeout    time.Duration `yaml:"read_timeout"`
	WriteTimeout   time.Duration `yaml:"write_timeout"`
	RequestTimeout time.Duration `yaml:"request_timeout"` // Таймаут на обработку запроса
}

// DatabaseConfig — конфигурация базы данных.
type DatabaseConfig struct {
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	Database        string        `yaml:"database"`
	SSLMode         string        `yaml:"ssl_mode"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
}

// DSN возвращает строку подключения к PostgreSQL в формате URL для pgx.
func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Database, c.SSLMode,
	)
}

// GraphQLConfig — конфигурация GraphQL.
type GraphQLConfig struct {
	EnablePlayground    bool `yaml:"enable_playground"`
	EnableIntrospection bool `yaml:"enable_introspection"`
	QueryCacheSize      int  `yaml:"query_cache_size"`
}

// LogLevel определяет уровень логирования.
type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

// LogConfig — конфигурация логгера.
type LogConfig struct {
	Level  LogLevel `yaml:"level"`
	Format string   `yaml:"format"` // "json" или "text"
}

// DefaultConfig возвращает конфигурацию по умолчанию.
func DefaultConfig() Config {
	return Config{
		Server: ServerConfig{
			Host:           "0.0.0.0",
			Port:           8080,
			ReadTimeout:    15 * time.Second,
			WriteTimeout:   15 * time.Second,
			RequestTimeout: 30 * time.Second,
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
		Log: LogConfig{
			Level:  LogLevelInfo,
			Format: "text",
		},
	}
}

// Addr возвращает адрес сервера.
func (c ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
