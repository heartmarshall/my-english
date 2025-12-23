package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/heartmarshall/my-english/graph"
	"github.com/jackc/pgx/v5/pgxpool"
)

// App — главная структура приложения.
type App struct {
	config Config
	logger *slog.Logger
	deps   *Dependencies
	server *http.Server
}

// New создаёт новое приложение.
func New(cfg Config) (*App, error) {
	// Инициализация логгера
	logger := NewLogger(cfg.Log)
	slog.SetDefault(logger)

	// Подключение к базе данных
	pool, err := connectDB(cfg.Database, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Инициализация зависимостей
	deps := NewDependencies(pool)

	// HTTP сервер
	server := newHTTPServer(cfg, deps, logger)

	return &App{
		config: cfg,
		logger: logger,
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
		a.logger.Info("starting server",
			slog.String("addr", a.config.Server.Addr()),
			slog.String("graphql", fmt.Sprintf("http://%s/graphql", a.config.Server.Addr())),
		)
		if a.config.GraphQL.EnablePlayground {
			a.logger.Info("playground enabled",
				slog.String("url", fmt.Sprintf("http://%s/playground", a.config.Server.Addr())),
			)
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
		a.logger.Info("received shutdown signal", slog.String("signal", sig.String()))
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
	a.deps.DB.Close()

	a.logger.Info("application stopped gracefully")
	return nil
}

// connectDB устанавливает соединение с базой данных через pgxpool.
func connectDB(cfg DatabaseConfig, logger *slog.Logger) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Настройка пула соединений
	config.MaxConns = int32(cfg.MaxOpenConns)
	config.MinConns = int32(cfg.MaxIdleConns)
	config.MaxConnLifetime = cfg.ConnMaxLifetime

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Проверка соединения
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("connected to database",
		slog.String("host", cfg.Host),
		slog.Int("port", cfg.Port),
		slog.String("database", cfg.Database),
		slog.Int("max_conns", cfg.MaxOpenConns),
	)
	return pool, nil
}

// newHTTPServer создаёт HTTP сервер.
func newHTTPServer(cfg Config, deps *Dependencies, logger *slog.Logger) *http.Server {
	mux := http.NewServeMux()

	// GraphQL handler
	graphqlServer := graph.NewServer(deps.Resolver, graph.ServerConfig{
		EnablePlayground:    cfg.GraphQL.EnablePlayground,
		EnableIntrospection: cfg.GraphQL.EnableIntrospection,
		QueryCacheSize:      cfg.GraphQL.QueryCacheSize,
	})
	graphqlServer.Routes(mux, "/graphql")

	// Health checks
	healthChecker := NewHealthChecker(deps.DB)
	mux.HandleFunc("/health", healthChecker.Handler())
	mux.HandleFunc("/live", healthChecker.LivenessHandler())
	mux.HandleFunc("/ready", healthChecker.ReadinessHandler())

	// DataLoader dependencies — использует сервис, не репозитории
	loaderDeps := graph.LoaderDeps{
		Loader: deps.Services.Loader,
	}

	// Middleware chain (порядок: снаружи → внутрь)
	// Recovery → Logging → Timeout → DataLoader → CORS → Handler
	handler := RecoveryMiddleware(logger)(
		LoggingMiddleware(logger)(
			TimeoutMiddleware(cfg.Server.RequestTimeout)(
				graph.DataLoaderMiddleware(loaderDeps)(
					graph.CORSMiddleware(mux),
				),
			),
		),
	)

	return &http.Server{
		Addr:         cfg.Server.Addr(),
		Handler:      handler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}
}
