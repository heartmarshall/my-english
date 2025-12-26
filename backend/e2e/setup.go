package e2e

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver for database/sql (pgx)
	_ "github.com/lib/pq"              // PostgreSQL driver for database/sql (lib/pq)
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestDB содержит информацию о тестовой базе данных.
type TestDB struct {
	Container testcontainers.Container
	Pool      *pgxpool.Pool
	DSN       string
}

// SetupTestDB создаёт тестовую базу данных с помощью testcontainers.
func SetupTestDB(ctx context.Context, t *testing.T) *TestDB {
	t.Helper()

	// Запускаем PostgreSQL контейнер
	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16-alpine"),
		postgres.WithDatabase("test_db"),
		postgres.WithUsername("test_user"),
		postgres.WithPassword("test_password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("Failed to start postgres container: %v", err)
	}

	// Получаем connection string
	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to get connection string: %v", err)
	}

	// Создаём пул соединений
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatalf("Failed to create connection pool: %v", err)
	}

	// Применяем миграции
	if err := runMigrations(ctx, connStr); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return &TestDB{
		Container: postgresContainer,
		Pool:      pool,
		DSN:       connStr,
	}
}

// Cleanup закрывает соединения и останавливает контейнер.
func (tdb *TestDB) Cleanup(ctx context.Context) error {
	if tdb.Pool != nil {
		tdb.Pool.Close()
	}
	if tdb.Container != nil {
		return tdb.Container.Terminate(ctx)
	}
	return nil
}

// runMigrations применяет миграции к тестовой базе данных.
func runMigrations(ctx context.Context, dsn string) error {
	// Получаем путь к директории с миграциями
	// Предполагаем, что мы запускаем тесты из корня проекта
	migrationsDir := "migrations"
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		// Пытаемся найти миграции относительно текущей директории
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}
		// Если мы в e2e директории, поднимаемся на уровень выше
		if filepath.Base(wd) == "e2e" {
			migrationsDir = filepath.Join("..", "migrations")
		}
	}

	// Используем goose API для применения миграций
	return runGooseMigrations(ctx, migrationsDir, dsn)
}

// runGooseMigrations применяет миграции через goose API.
func runGooseMigrations(ctx context.Context, migrationsDir, dsn string) error {
	absMigrationsDir, err := filepath.Abs(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path to migrations: %w", err)
	}

	// Используем goose API напрямую
	// goose.OpenDBWithDriver использует database/sql, поэтому нужен драйвер lib/pq
	db, err := goose.OpenDBWithDriver("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Проверяем соединение
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Применяем миграции
	if err := goose.UpContext(ctx, db, absMigrationsDir); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
