# Инструкция по работе с API через Postman

## Подключение к сервису в Docker

### 1. Проверка, что сервис запущен

```bash
docker compose ps
```

Убедитесь, что контейнер `my-english-backend` работает.

### 2. Настройка Postman

#### Базовые параметры запроса:

- **Метод**: `POST`
- **URL**: `http://localhost:8080/graphql`
- **Headers**:
  - `Content-Type: application/json`

#### Body (raw JSON):

```json
{
  "query": "query Suggest($query: String!) { suggest(query: $query) { text transcription translations definition origin existingWordId } }",
  "variables": {
    "query": "hello"
  }
}
```

### 3. Примеры запросов

#### Полный запрос для автоподсказок (все поля):

```json
{
  "query": "query Suggest($query: String!) { suggest(query: $query) { text transcription translations definition origin existingWordId } }",
  "variables": {
    "query": "hello"
  }
}
```

#### Минимальный запрос (только текст и транскрипция):

```json
{
  "query": "query Suggest($query: String!) { suggest(query: $query) { text transcription } }",
  "variables": {
    "query": "hello"
  }
}
```

#### Запрос с переводом и определением:

```json
{
  "query": "query Suggest($query: String!) { suggest(query: $query) { text transcription translations definition } }",
  "variables": {
    "query": "hello"
  }
}
```

### 4. Почему возвращается только транскрипция?

**Проблема**: В GraphQL нужно явно указывать все поля, которые вы хотите получить в ответе.

Если ваш запрос выглядит так:
```graphql
query {
  suggest(query: "hello") {
    transcription
  }
}
```

То вы получите только `transcription`, даже если в базе данных есть и другие поля.

**Решение**: Добавьте все нужные поля в запрос:
```graphql
query {
  suggest(query: "hello") {
    text
    transcription
    translations
    definition
    origin
    existingWordId
  }
}
```

### 5. Альтернатива: GraphQL Playground

Если Playground включен (по умолчанию включен), вы можете использовать браузерный интерфейс:

- **URL**: `http://localhost:8080/playground`

Это удобнее для тестирования, так как там есть автодополнение и документация схемы.

### 6. Проверка доступности сервиса

Перед запросами к GraphQL можно проверить health endpoint:

- **Метод**: `GET`
- **URL**: `http://localhost:8080/health`

Ожидаемый ответ: `200 OK` с JSON объектом статуса.



