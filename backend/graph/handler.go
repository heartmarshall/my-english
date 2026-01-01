package graph

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	ctxlog "github.com/heartmarshall/my-english/pkg/context"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// ServerConfig — конфигурация GraphQL сервера.
type ServerConfig struct {
	// EnablePlayground включает GraphQL Playground.
	EnablePlayground bool
	// EnableIntrospection включает introspection запросы.
	EnableIntrospection bool
	// QueryCacheSize — размер кэша для запросов.
	QueryCacheSize int
}

// DefaultServerConfig возвращает конфигурацию по умолчанию.
func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		EnablePlayground:    true,
		EnableIntrospection: true,
		QueryCacheSize:      1000,
	}
}

// NewHandler создаёт HTTP handler для GraphQL.
func NewHandler(resolver *Resolver, cfg ServerConfig) http.Handler {
	srv := handler.New(NewExecutableSchema(Config{
		Resolvers: resolver,
	}))

	// Логирование ошибок GraphQL
	srv.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		logger := ctxlog.L(ctx)
		logger.Error("graphql panic recovered",
			slog.Any("error", err),
		)
		return gqlerror.Errorf("internal server error")
	})

	srv.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
		logger := ctxlog.L(ctx)
		logger.Error("graphql error",
			slog.String("error", e.Error()),
		)
		// Возвращаем ошибку как есть, чтобы она попала в ответ
		if gqlErr, ok := e.(*gqlerror.Error); ok {
			return gqlErr
		}
		return gqlerror.Errorf(e.Error())
	})

	// Transports
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	// Кэширование
	srv.SetQueryCache(lru.New[*ast.QueryDocument](cfg.QueryCacheSize))

	// Расширения
	if cfg.EnableIntrospection {
		srv.Use(extension.Introspection{})
	}
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return srv
}

// NewPlaygroundHandler создаёт handler для GraphQL Playground.
func NewPlaygroundHandler(endpoint string) http.Handler {
	return playground.Handler("GraphQL Playground", endpoint)
}

// Server — HTTP сервер для GraphQL.
type Server struct {
	graphqlHandler http.Handler
	config         ServerConfig
}

// NewServer создаёт новый GraphQL сервер.
func NewServer(resolver *Resolver, cfg ServerConfig) *Server {
	return &Server{
		graphqlHandler: NewHandler(resolver, cfg),
		config:         cfg,
	}
}

// Handler возвращает HTTP handler для GraphQL.
func (s *Server) Handler() http.Handler {
	return s.graphqlHandler
}

// Routes регистрирует маршруты на mux.
func (s *Server) Routes(mux *http.ServeMux, path string) {
	mux.Handle(path, s.graphqlHandler)

	if s.config.EnablePlayground {
		mux.Handle("/playground", NewPlaygroundHandler(path))
	}
}

// --- Middleware ---

// CORSMiddleware добавляет CORS заголовки.
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RecoveryMiddleware обрабатывает паники.
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// TODO: добавить логирование
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware логирует HTTP запросы.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap response writer to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		// TODO: заменить на структурированный логгер
		_ = time.Since(start)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
