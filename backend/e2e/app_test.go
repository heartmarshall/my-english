package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/heartmarshall/my-english/internal/app"
	"github.com/heartmarshall/my-english/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TestApp содержит тестовое приложение и его зависимости.
type TestApp struct {
	App      *app.App
	Server   *http.Server
	DB       *pgxpool.Pool
	URL      string
	listener net.Listener
}

// SetupTestApp создаёт тестовое приложение с подключением к тестовой БД.
func SetupTestApp(ctx context.Context, t *testing.T, db *pgxpool.Pool) *TestApp {
	t.Helper()

	// Создаём конфигурацию для тестов
	cfg := config.Config{
		Server: config.ServerConfig{
			Host:           "127.0.0.1",
			Port:           0, // 0 означает автоматический выбор порта через net.Listen
			ReadTimeout:    15 * time.Second,
			WriteTimeout:   15 * time.Second,
			RequestTimeout: 30 * time.Second,
		},
		Database: config.DatabaseConfig{
			Host:            "localhost",
			Port:            5432,
			User:            "test_user",
			Password:        "test_password",
			Database:        "test_db",
			SSLMode:         "disable",
			MaxOpenConns:    10,
			MaxIdleConns:    5,
			ConnMaxLifetime: 5 * time.Minute,
		},
		GraphQL: config.GraphQLConfig{
			EnablePlayground:    true,
			EnableIntrospection: true,
			QueryCacheSize:      1000,
		},
		Log: config.LogConfig{
			Level:  "error", // Минимальное логирование в тестах
			Format: "text",
		},
	}

	// Создаём приложение с готовым пулом БД
	application, err := app.NewWithPool(cfg, db)
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	// Запускаем сервер на случайном порту
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}

	// Получаем адрес сервера
	addr := listener.Addr().(*net.TCPAddr)
	url := fmt.Sprintf("http://127.0.0.1:%d", addr.Port)

	// Обновляем адрес сервера
	server := application.Server()
	server.Addr = fmt.Sprintf(":%d", addr.Port)

	// Запускаем сервер в горутине
	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			t.Logf("Server error: %v", err)
		}
	}()

	// Ждём, пока сервер запустится
	time.Sleep(200 * time.Millisecond)

	return &TestApp{
		App:      application,
		Server:   server,
		DB:       db,
		URL:      url,
		listener: listener,
	}
}

// GraphQLRequest представляет GraphQL запрос.
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// GraphQLResponse представляет GraphQL ответ.
type GraphQLResponse struct {
	Data   json.RawMessage        `json:"data"`
	Errors []GraphQLError         `json:"errors,omitempty"`
}

// GraphQLError представляет ошибку GraphQL.
type GraphQLError struct {
	Message    string                 `json:"message"`
	Path       []interface{}          `json:"path,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

// DoGraphQLRequest выполняет GraphQL запрос.
func (ta *TestApp) DoGraphQLRequest(ctx context.Context, req GraphQLRequest) (*GraphQLResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", ta.URL+"/graphql", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(respBody))
	}

	var graphqlResp GraphQLResponse
	if err := json.Unmarshal(respBody, &graphqlResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &graphqlResp, nil
}

// Cleanup останавливает тестовое приложение.
func (ta *TestApp) Cleanup(ctx context.Context) error {
	if ta.listener != nil {
		ta.listener.Close()
	}
	if ta.App != nil {
		return ta.App.Shutdown()
	}
	return nil
}

