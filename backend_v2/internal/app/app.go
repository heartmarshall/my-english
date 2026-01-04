package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/heartmarshall/my-english/graph"
	"github.com/heartmarshall/my-english/internal/clients/freedict"
	"github.com/heartmarshall/my-english/internal/config"
	"github.com/heartmarshall/my-english/internal/database"
	repository "github.com/heartmarshall/my-english/internal/database/repository/factory"
	"github.com/heartmarshall/my-english/internal/service/card"
	"github.com/heartmarshall/my-english/internal/service/dictionary"
	"github.com/heartmarshall/my-english/internal/service/inbox"
	"github.com/heartmarshall/my-english/internal/service/loader"
	"github.com/heartmarshall/my-english/internal/service/stats"
	"github.com/heartmarshall/my-english/internal/service/study"
	"github.com/heartmarshall/my-english/internal/transport/middleware"
)

// App представляет основное приложение со всеми зависимостями.
type App struct {
	cfg       config.Config
	logger    *slog.Logger
	pool      *pgxpool.Pool
	txManager *database.TxManager
	repos     *repository.Factory
	services  *Services
	resolver  *graph.Resolver
	server    *http.Server
	health    *HealthChecker
}

// Services содержит все сервисы приложения.
type Services struct {
	Card       *card.Service
	Dictionary *dictionary.Service
	Inbox      *inbox.Service
	Loader     *loader.Service
	Stats      *stats.Service
	Study      *study.Service
}

// New создаёт и инициализирует новое приложение.
func New(cfg config.Config) (*App, error) {
	// 1. Инициализация логгера
	logger := NewLogger(cfg.Log)

	// 2. Подключение к базе данных
	pool, err := initDatabase(cfg.Database, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// 3. Инициализация менеджера транзакций
	txManager := database.NewTxManager(pool)

	// 4. Инициализация фабрики репозиториев
	repos := repository.NewFactory()

	// 5. Инициализация сервисов
	services, err := initServices(repos, txManager, pool, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize services: %w", err)
	}

	// 6. Инициализация GraphQL resolver
	resolver := &graph.Resolver{
		CardService:       services.Card,
		DictionaryService: services.Dictionary,
		InboxService:      services.Inbox,
		LoaderService:     services.Loader,
		StatsService:      services.Stats,
		StudyService:      services.Study,
	}

	// 7. Инициализация health checker
	health := NewHealthChecker(pool)

	// 8. Создание HTTP сервера
	server := initHTTPServer(cfg, resolver, services.Loader, health, logger)

	app := &App{
		cfg:       cfg,
		logger:    logger,
		pool:      pool,
		txManager: txManager,
		repos:     repos,
		services:  services,
		resolver:  resolver,
		server:    server,
		health:    health,
	}

	return app, nil
}

// Run запускает HTTP сервер.
func (a *App) Run() error {
	a.logger.Info("starting server",
		slog.String("address", a.cfg.Server.Addr()),
	)

	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}

// Shutdown корректно останавливает приложение.
func (a *App) Shutdown(ctx context.Context) error {
	a.logger.Info("shutting down application")

	// Закрываем HTTP сервер
	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Error("failed to shutdown server", slog.Any("error", err))
	}

	// Закрываем пул соединений
	a.pool.Close()

	a.logger.Info("application shutdown complete")
	return nil
}

// initDatabase инициализирует подключение к базе данных.
func initDatabase(cfg config.DatabaseConfig, logger *slog.Logger) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}

	config.MaxConns = int32(cfg.MaxOpenConns)
	config.MaxConnIdleTime = cfg.ConnMaxLifetime
	config.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Проверяем соединение
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("database connection established",
		slog.String("host", cfg.Host),
		slog.Int("port", cfg.Port),
		slog.String("database", cfg.Database),
	)

	return pool, nil
}

// initServices инициализирует все сервисы приложения.
func initServices(
	repos *repository.Factory,
	txManager *database.TxManager,
	pool *pgxpool.Pool,
	logger *slog.Logger,
) (*Services, error) {
	// Инициализация провайдеров словаря
	providers := []dictionary.Provider{
		freedict.NewClient(),
	}

	// Card Service
	cardService := card.New(card.Deps{
		Repos:     repos,
		TxManager: txManager,
	})

	// Dictionary Service
	dictionaryService := dictionary.New(dictionary.Deps{
		Repos:     repos,
		TxManager: txManager,
		Providers: providers,
	})

	// Inbox Service
	inboxService := inbox.New(inbox.Deps{
		Repos:     repos,
		TxManager: txManager,
	})

	// Loader Service (для DataLoaders)
	loaderService := loader.New(repos, pool)

	// Stats Service
	statsService := stats.New(stats.Deps{
		Repos:     repos,
		TxManager: txManager,
	})

	// Study Service
	studyService := study.New(study.Deps{
		Repos:     repos,
		TxManager: txManager,
		// Algorithm будет создан автоматически (SM2 по умолчанию)
	})

	return &Services{
		Card:       cardService,
		Dictionary: dictionaryService,
		Inbox:      inboxService,
		Loader:     loaderService,
		Stats:      statsService,
		Study:      studyService,
	}, nil
}

// initHTTPServer создаёт и настраивает HTTP сервер.
func initHTTPServer(
	cfg config.Config,
	resolver *graph.Resolver,
	loaderService *loader.Service,
	health *HealthChecker,
	logger *slog.Logger,
) *http.Server {
	mux := http.NewServeMux()

	// GraphQL handler
	graphqlHandler := createGraphQLHandler(cfg.GraphQL, resolver, loaderService)
	mux.Handle("/graphql", graphqlHandler)

	// GraphQL Playground (если включен)
	if cfg.GraphQL.EnablePlayground {
		playgroundHandler := playground.Handler("GraphQL Playground", "/graphql")
		mux.Handle("/playground", playgroundHandler)
		logger.Info("GraphQL playground enabled at /playground")
	}

	// Health check endpoints
	mux.HandleFunc("/health", health.Handler())
	mux.HandleFunc("/health/live", health.LivenessHandler())
	mux.HandleFunc("/health/ready", health.ReadinessHandler())

	// Применяем middleware
	// Порядок важен: сначала DataLoader (чтобы создать loaders в контексте),
	// затем остальные middleware
	var httpHandler http.Handler = mux
	httpHandler = graph.Middleware(loaderService)(httpHandler)
	httpHandler = middleware.LoggingMiddleware(logger)(httpHandler)
	httpHandler = middleware.RecoveryMiddleware(logger)(httpHandler)
	httpHandler = middleware.TimeoutMiddleware(cfg.Server.RequestTimeout)(httpHandler)

	return &http.Server{
		Addr:         cfg.Server.Addr(),
		Handler:      httpHandler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}
}

// createGraphQLHandler создаёт GraphQL handler с настройками.
func createGraphQLHandler(
	cfg config.GraphQLConfig,
	resolver *graph.Resolver,
	loaderService *loader.Service,
) http.Handler {
	// Создаём executable schema
	// Schema, Directives и Complexity будут использованы из generated.go по умолчанию
	execSchema := graph.NewExecutableSchema(graph.Config{
		Resolvers: resolver,
		// Schema будет nil, что означает использование parsedSchema из generated.go
		// Directives и Complexity будут пустыми структурами по умолчанию
	})

	// Создаём GraphQL handler
	h := handler.NewDefaultServer(execSchema)

	// Настройка кэширования запросов
	if cfg.QueryCacheSize > 0 {
		h.Use(extension.FixedComplexityLimit(1000)) // Защита от сложных запросов
	}

	// Отключение introspection, если нужно
	if !cfg.EnableIntrospection {
		h.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
			opCtx := graphql.GetOperationContext(ctx)
			if opCtx.OperationName == "IntrospectionQuery" {
				return graphql.OneShot(graphql.ErrorResponse(ctx, "introspection is disabled"))
			}
			return next(ctx)
		})
	}

	return h
}
