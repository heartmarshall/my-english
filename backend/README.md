# my-english

Стек: Go + postgres + graphql + squirel (для бд)

Сервис для практики английского языка.

## Генерация GraphQL кода

Для генерации GraphQL кода из схемы используйте команду:

```bash
make generate
```

Или напрямую:

```bash
go run github.com/99designs/gqlgen generate
```

После генерации будут созданы файлы:
- `graph/generated.go` - сгенерированный код сервера
- `graph/models_gen.go` - сгенерированные модели
- `graph/resolvers.go` - файлы резолверов (создаются автоматически)
