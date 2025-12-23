package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/heartmarshall/my-english/graph"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// App — главная структура приложения.
type App struct {
	config Config
	deps   *Dependencies
	server *http.Server
}

// New создаёт новое приложение.
func New(cfg Config) (*App, error) {
	// Подключение к базе данных
	db, err := connectDB(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Инициализация зависимостей
	deps := NewDependencies(db)

	// HTTP сервер
	server := newHTTPServer(cfg, deps.Resolver)

	return &App{
		config: cfg,
		deps:   deps,
		server: server,
	}, nil
}

// Run запускает приложение и ожидает сигнала завершения.
func (a *App) Run() error {
	// Канал для ошибок сервера
	errChan := make(chan error, 1)

	// Запуск сервера в горутине
	go func() {
		log.Printf("Starting server on %s", a.config.Server.Addr())
		log.Printf("GraphQL endpoint: http://%s/graphql", a.config.Server.Addr())
		if a.config.GraphQL.EnablePlayground {
			log.Printf("GraphQL Playground: http://%s/playground", a.config.Server.Addr())
		}

		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Ожидание сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		return fmt.Errorf("server error: %w", err)
	case sig := <-quit:
		log.Printf("Received signal: %v. Shutting down...", sig)
	}

	return a.Shutdown()
}

// Shutdown gracefully завершает приложение.
func (a *App) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Остановка HTTP сервера
	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	// Закрытие соединения с БД
	if err := a.deps.DB.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}

	log.Println("Application stopped gracefully")
	return nil
}

// connectDB устанавливает соединение с базой данных.
func connectDB(cfg DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, err
	}

	// Настройка пула соединений
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Проверка соединения
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Connected to database")
	return db, nil
}

// newHTTPServer создаёт HTTP сервер.
func newHTTPServer(cfg Config, resolver *graph.Resolver) *http.Server {
	mux := http.NewServeMux()

	// GraphQL handler
	graphqlServer := graph.NewServer(resolver, graph.ServerConfig{
		EnablePlayground:    cfg.GraphQL.EnablePlayground,
		EnableIntrospection: cfg.GraphQL.EnableIntrospection,
		QueryCacheSize:      cfg.GraphQL.QueryCacheSize,
	})
	graphqlServer.Routes(mux, "/graphql")

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Middleware chain
	handler := graph.RecoveryMiddleware(
		graph.LoggingMiddleware(
			graph.CORSMiddleware(mux),
		),
	)

	return &http.Server{
		Addr:         cfg.Server.Addr(),
		Handler:      handler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}
}
