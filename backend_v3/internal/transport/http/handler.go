package http

import (
	"log/slog"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/heartmarshall/my-english/graph"
	"github.com/heartmarshall/my-english/internal/config"
	"github.com/heartmarshall/my-english/internal/database/repository"
	"github.com/heartmarshall/my-english/internal/service"
	"github.com/heartmarshall/my-english/internal/transport/dataloader"
	"github.com/heartmarshall/my-english/internal/transport/middleware"
	"github.com/vektah/gqlparser/v2/ast"
)

// NewHandler создает и настраивает основной HTTP хендлер приложения.
func NewHandler(
	cfg *config.Config,
	logger *slog.Logger,
	services *service.Services,
	repos *repository.Registry,
) http.Handler {
	// 1. Инициализация GraphQL схемы
	gqlConfig := graph.Config{
		Resolvers: &graph.Resolver{
			Services: services,
		},
	}

	// Создаем сервер gqlgen
	srv := handler.New(graph.NewExecutableSchema(gqlConfig))

	// 2. Настройка транспортов GraphQL
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	// 3. Настройка кэширования запросов и сложности
	srv.SetQueryCache(lru.New[*ast.QueryDocument](cfg.GraphQL.QueryCacheSize))
	srv.Use(extension.Introspection{}) // Включаем интроспекцию (можно отключить для prod)
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	// 4. Настройка роутера (используем стандартный ServeMux)
	mux := http.NewServeMux()

	// Health check
	// (Если health.go в пакете app, его можно подключить в main.go или передать сюда,
	// но для простоты опустим, так как он обычно отдельный)

	// GraphQL Playground (только если включен)
	if cfg.GraphQL.EnablePlayground {
		mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	}

	// Основной Endpoint
	mux.Handle("/query", srv)

	// 5. Подключение Middleware (порядок важен!)
	// Запросы проходят снизу вверх (в коде) или снаружи внутрь.

	// Базовый хендлер
	var handler http.Handler = mux

	// Внедряем DataLoaders (обязательно перед GraphQL хендлером)
	handler = dataloader.Middleware(repos)(handler)

	// Логирование, Recovery, Timeout
	handler = middleware.TimeoutMiddleware(cfg.Server.RequestTimeout)(handler)
	handler = middleware.LoggingMiddleware(logger)(handler)
	handler = middleware.RecoveryMiddleware(logger)(handler)

	return handler
}
