# Быстрый старт с Docker

## Первый запуск

1. **Скопируйте файл с переменными окружения** (опционально, можно использовать значения по умолчанию):
   ```bash
   cp .docker-compose.env.example .env
   ```

2. **Запустите все сервисы**:
   ```bash
   docker compose up -d
   ```
   
   Это запустит:
   - PostgreSQL базу данных
   - Миграции (автоматически)
   - Backend сервис

3. **Проверьте статус**:
   ```bash
   docker compose ps
   ```

4. **Проверьте логи** (если нужно):
   ```bash
   docker compose logs -f backend
   ```

## Доступ к сервисам

- **GraphQL API**: http://localhost:8080/graphql
- **GraphQL Playground**: http://localhost:8080/playground
- **Health Check**: http://localhost:8080/health
- **PostgreSQL**: localhost:5432

## Полезные команды

```bash
# Остановить все сервисы
docker compose down

# Перезапустить сервисы
docker compose restart

# Просмотр логов
docker compose logs -f [service_name]

# Подключиться к базе данных
docker compose exec postgres psql -U postgres -d my_english

# Запустить миграции вручную
docker compose up migrate
```

## Troubleshooting

Если что-то не работает:

1. Проверьте, что все контейнеры запущены:
   ```bash
   docker compose ps
   ```

2. Проверьте логи:
   ```bash
   docker compose logs
   ```

3. Пересоберите образы:
   ```bash
   docker compose build --no-cache
   docker compose up -d
   ```

4. Очистите все и начните заново:
   ```bash
   docker compose down -v
   docker compose up -d
   ```

Подробная документация: [DOCKER.md](./DOCKER.md)

