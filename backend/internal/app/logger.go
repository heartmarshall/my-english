package app

import (
	"context"
	"log/slog"
	"os"
)

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
	Level  LogLevel
	Format string // "json" или "text"
}

// DefaultLogConfig возвращает конфигурацию по умолчанию.
func DefaultLogConfig() LogConfig {
	return LogConfig{
		Level:  LogLevelInfo,
		Format: "text",
	}
}

// NewLogger создаёт новый логгер.
func NewLogger(cfg LogConfig) *slog.Logger {
	var level slog.Level
	switch cfg.Level {
	case LogLevelDebug:
		level = slog.LevelDebug
	case LogLevelWarn:
		level = slog.LevelWarn
	case LogLevelError:
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	if cfg.Format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}

// --- Context helpers ---

type ctxKey string

const loggerKey ctxKey = "logger"

// WithLogger добавляет логгер в контекст.
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// LoggerFromContext извлекает логгер из контекста.
// Если логгер не найден, возвращает default логгер.
func LoggerFromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}

// L — короткий алиас для LoggerFromContext.
func L(ctx context.Context) *slog.Logger {
	return LoggerFromContext(ctx)
}
