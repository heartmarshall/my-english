# E2E Tests

End-to-end тесты для сервиса my-english с использованием testcontainers.

## Требования

- Docker (для запуска testcontainers)
- Go 1.24+

## Запуск тестов

```bash
# Запустить все e2e тесты
go test ./e2e/... -v

# Запустить конкретный тест
go test ./e2e/... -run TestCreateWord -v

# Запустить с таймаутом (для CI/CD)
go test ./e2e/... -timeout 5m -v
```

## Структура

- `setup.go` - утилиты для настройки тестовой БД с testcontainers
- `app_test.go` - утилиты для создания тестового приложения
- `words_test.go` - тесты для операций со словами
- `inbox_test.go` - тесты для операций с inbox

## Как это работает

1. **SetupTestDB** - создаёт PostgreSQL контейнер через testcontainers
2. **runMigrations** - применяет миграции к тестовой БД
3. **SetupTestApp** - создаёт тестовое приложение с подключением к тестовой БД
4. Тесты выполняют GraphQL запросы через HTTP

## Добавление новых тестов

1. Создайте новый файл `*_test.go` в директории `e2e/`
2. Используйте `SetupTestDB` и `SetupTestApp` для настройки
3. Используйте `DoGraphQLRequest` для выполнения GraphQL запросов
4. Не забудьте вызвать `Cleanup` в `defer`

## Пример

```go
func TestMyFeature(t *testing.T) {
    ctx := context.Background()
    
    testDB := SetupTestDB(ctx, t)
    defer testDB.Cleanup(ctx)
    
    testApp := SetupTestApp(ctx, t, testDB.Pool)
    defer testApp.Cleanup(ctx)
    
    // Ваш тест...
}
```

