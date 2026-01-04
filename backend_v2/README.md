# my-english

Стек: Go + postgres + graphql + squirel (для бд)

Сервис для практики английского языка.

## Быстрый старт с Docker

Самый простой способ запустить проект - использовать Docker:

```bash
# 1. Скопируйте файл с переменными окружения
cp .docker-compose.env.example .env

# 2. Запустите все сервисы
docker compose up -d

# 3. Проверьте статус
docker compose ps
```

Подробная документация по Docker: [DOCKER.md](./DOCKER.md)

## Локальная разработка

### Требования

- Go 1.24+
- PostgreSQL 16+
- Make

### Установка

```bash
# Установить зависимости
go mod download

# Сгенерировать GraphQL код
make generate
```

### Запуск

```bash
# Запустить сервер
go run cmd/server/main.go

# Или с конфигом
go run cmd/server/main.go -config config.example.yaml
```

### Миграции

```bash
# Применить миграции
make migrate-up DB_URL=postgres://user:password@localhost:5432/my_english?sslmode=disable

# Откатить последнюю миграцию
make migrate-down DB_URL=postgres://user:password@localhost:5432/my_english?sslmode=disable

# Статус миграций
make migrate-status DB_URL=postgres://user:password@localhost:5432/my_english?sslmode=disable
```

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

## Docker команды

Основные команды для работы с Docker:

```bash
# Запустить все сервисы
docker compose up -d

# Запустить миграции
docker compose up migrate

# Просмотр логов
docker compose logs -f [service_name]

# Остановить все
docker compose down
```

Подробная документация: [DOCKER.md](./DOCKER.md)


# Docker Setup

## Быстрый старт

```bash
# 1. Опционально: скопировать переменные окружения
cp .docker-compose.env.example .env

# 2. Запустить все сервисы
docker compose up -d

# 3. Проверить статус
docker compose ps
```

## Основные команды

```bash
# Запустить сервисы
docker compose up -d

# Остановить сервисы
docker compose down

# Просмотр логов
docker compose logs -f [service_name]

# Перезапустить
docker compose restart

# Запустить миграции вручную
docker compose up migrate

# Подключиться к БД
docker compose exec postgres psql -U postgres -d my_english
```

## Доступ к сервисам

- **GraphQL API**: http://localhost:8080/graphql
- **GraphQL Playground**: http://localhost:8080/playground
- **Health Check**: http://localhost:8080/health
- **PostgreSQL**: localhost:5432

## Переменные окружения

Основные переменные (можно настроить в `.env`):

- `DB_USER`, `DB_PASSWORD`, `DB_NAME` - настройки БД
- `SERVER_PORT` - порт backend (по умолчанию: 8080)
- `LOG_LEVEL` - уровень логирования

Полный список: `.docker-compose.env.example`

## Важные замечания

### Использование Docker Compose v2

Используйте `docker compose` (v2) вместо `docker-compose` (v1):

```bash
docker compose up -d  # правильно
docker-compose up -d  # устаревший способ
```

### Проблемы с сетями Docker

Если ошибка `all predefined address pools have been fully subnetted`:

```bash
docker network prune -f
docker system prune -a --volumes -f
```

## Troubleshooting

**Backend не подключается к БД**: Проверьте, что используете `postgres` как hostname (не `localhost`)

**Миграции не применяются**: 
```bash
docker compose up migrate
```

**Порт занят**: Измените порт в `.env` или `docker-compose.yml`
