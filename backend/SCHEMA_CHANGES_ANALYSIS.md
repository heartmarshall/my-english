# Анализ изменений GraphQL схемы

## 📋 Обзор изменений

### 1. **Изменение типа `Meaning.translationRu`**
**Было:** `translationRu: String!`  
**Стало:** `translationRu: [string]` (массив строк, nullable)

**Проблема:** В базе данных `translation_ru` хранится как одно поле `TEXT NOT NULL`, но в GraphQL схеме теперь это массив.

**Решения:**
- **Вариант A (временный):** Возвращать массив из одного элемента `[m.TranslationRu]`
- **Вариант B (правильный):** Создать отдельную таблицу `translations` для хранения множественных переводов
- **Вариант C:** Хранить JSON массив в поле `translation_ru`

**Что нужно обновить:**
- `backend/graph/converter.go` - функция `ToGraphQLMeaning()` и `ToGraphQLMeaningBasic()` должны возвращать `[]string` вместо `string`
- Добавить field resolver для `Meaning.translationRu` в `backend/graph/schema.resolvers.go`
- Обновить модель `model.Meaning` (если выбран вариант B или C)

---

### 2. **Переименование `AddWordInput` → `CreateWordInput`**
**Было:** `input AddWordInput`  
**Стало:** `input CreateWordInput`

**Что нужно обновить:**
- `backend/graph/schema.resolvers.go`:
  - `CreateWord(ctx, input AddWordInput)` → `CreateWord(ctx, input CreateWordInput)`
  - `UpdateWord(ctx, id string, input AddWordInput)` → `UpdateWord(ctx, id string, input CreateWordInput)`
- `backend/graph/converter.go`:
  - `ToCreateWordInput(input AddWordInput)` → `ToCreateWordInput(input CreateWordInput)`
  - `ToUpdateWordInput(input AddWordInput)` → `ToUpdateWordInput(input CreateWordInput)`
- После `make generate` все ссылки на `AddWordInput` в `generated.go` и `models_gen.go` обновятся автоматически

---

### 3. **Изменение `words` query - пагинация через Connection**
**Было:** 
```graphql
words(filter: WordFilter, limit: Int, offset: Int): [Word!]!
```

**Стало:**
```graphql
words(filter: WordFilter, first: Int = 20, after: String): WordConnection!
```

**Что нужно обновить:**
- `backend/graph/schema.resolvers.go`:
  - Изменить сигнатуру `Words()`:
    ```go
    // Было:
    Words(ctx context.Context, filter *WordFilter, limit *int, offset *int) ([]*Word, error)
    
    // Стало:
    Words(ctx context.Context, filter *WordFilter, first *int, after *string) (*WordConnection, error)
    ```
- Реализовать cursor-based пагинацию:
  - Конвертировать `after` cursor в offset (или использовать cursor напрямую)
  - Вычислить `hasNextPage` и `endCursor`
  - Вернуть `WordConnection` с `edges`, `pageInfo`, `totalCount`
- `backend/graph/converter.go`:
  - Добавить функции для создания `WordConnection`, `WordEdge`, `PageInfo`
- `backend/internal/service/word/`:
  - Возможно, добавить метод для получения `totalCount` с фильтром

---

### 4. **Новые типы (отсутствуют в схеме, но используются)**
**Проблема:** В мутациях используются `CreateWordPayload` и `UpdateWordPayload`, но они не определены в схеме!

**Нужно добавить в `schema.graphqls`:**
```graphql
type CreateWordPayload {
  word: Word!
}

type UpdateWordPayload {
  word: Word!
}
```

**Что нужно обновить:**
- Добавить определения типов в `schema.graphqls`
- После `make generate` появятся в `models_gen.go`
- Обновить резолверы `CreateWord` и `UpdateWord` для возврата payload вместо `Word`

---

### 5. **Новые типы: `InboxItem`, `Suggestion`**
**Новые типы:**
- `InboxItem` - элемент корзины входящих слов
- `Suggestion` - подсказка для автокомплита
- `SuggestionOrigin` enum (LOCAL, DICTIONARY)

**Что нужно реализовать:**
- Создать таблицу `inbox_items` в БД (миграция)
- Создать модель `model.InboxItem`
- Создать репозиторий `internal/database/inbox/`
- Создать сервис `internal/service/inbox/` (или добавить в существующий)
- Добавить резолверы в `schema.resolvers.go`:
  - `InboxItems(ctx) ([]*InboxItem, error)`
  - `Suggest(ctx, query string) ([]*Suggestion, error)`
  - `AddToInbox(ctx, text string, sourceContext *string) (*InboxItem, error)`
  - `DeleteInboxItem(ctx, id string) (bool, error)`
  - `ConvertInboxItem(ctx, inboxId string, input CreateWordInput) (*CreateWordPayload, error)`
- Добавить конвертеры в `converter.go`

---

### 6. **Обновление `WordFilter`**
**Было:**
```graphql
input WordFilter {
  search: String
}
```

**Стало:**
```graphql
input WordFilter {
  search: String
  status: LearningStatus
  tags: [String!]
}
```

**Что нужно обновить:**
- `backend/graph/models_gen.go` - после `make generate` обновится автоматически
- `backend/graph/converter.go` - функция `ToWordFilter()` должна обрабатывать новые поля
- `backend/internal/service/word/dto.go` - обновить `WordFilter` структуру
- `backend/internal/service/word/list.go` - добавить фильтрацию по `status` и `tags`

---

### 7. **Новое поле `sourceContext` в `CreateWordInput`**
**Добавлено:** `sourceContext: String` в `CreateWordInput`

**Что нужно обновить:**
- `backend/graph/models_gen.go` - после `make generate` появится поле
- `backend/graph/converter.go` - `ToCreateWordInput()` должен обрабатывать `sourceContext`
- `backend/internal/service/word/dto.go` - добавить `SourceContext *string` в `CreateWordInput`
- Решить, где хранить `sourceContext`:
  - В таблице `words` (добавить колонку `source_context`)
  - В отдельной таблице для истории источников

---

## 📝 План действий (приоритетный порядок)

### Этап 1: Критические изменения (без них код не скомпилируется)
1. ✅ Добавить определения `CreateWordPayload` и `UpdateWordPayload` в схему
2. ✅ Перегенерировать GraphQL код (`make generate`)
3. ✅ Обновить резолверы для использования `CreateWordInput` вместо `AddWordInput`
4. ✅ Обновить `words` query для возврата `WordConnection`
5. ✅ Добавить field resolver для `Meaning.translationRu` (массив)

### Этап 2: Обновление конвертеров и логики
6. ✅ Обновить `ToGraphQLMeaning()` для работы с массивом `translationRu`
7. ✅ Обновить `ToWordFilter()` для обработки `status` и `tags`
8. ✅ Обновить `ToCreateWordInput()` для обработки `sourceContext`
9. ✅ Реализовать cursor-based пагинацию

### Этап 3: Новый функционал (Inbox, Suggestions)
10. ⏳ Создать миграцию для таблицы `inbox_items`
11. ⏳ Создать модель `InboxItem`
12. ⏳ Создать репозиторий и сервис для Inbox
13. ⏳ Реализовать резолверы для Inbox операций
14. ⏳ Реализовать резолвер `suggest()` (поиск в словаре и внешних источниках)

### Этап 4: Решение проблемы `translationRu` как массива
15. ⏳ Выбрать стратегию (A, B или C)
16. ⏳ Реализовать выбранную стратегию
17. ⏳ Обновить миграции/модели при необходимости

---

## 🔍 Детали по каждому изменению

### Изменение 1: `translationRu: [string]`

**Текущая реализация:**
- БД: `translation_ru TEXT NOT NULL` (одно значение)
- Модель: `TranslationRu string`
- GraphQL: было `String!`, стало `[string]`

**Рекомендация:** Начать с варианта A (временный), затем перейти к варианту B (правильный).

**Код для варианта A:**
```go
// В converter.go
func ToGraphQLMeaningBasic(m *model.Meaning) *Meaning {
    // ...
    translationRuArray := []string{}
    if m.TranslationRu != "" {
        translationRuArray = []string{m.TranslationRu}
    }
    return &Meaning{
        // ...
        TranslationRu: translationRuArray,
        // ...
    }
}
```

**Код для варианта B (будущее):**
- Создать таблицу `translations`:
  ```sql
  CREATE TABLE translations (
      id SERIAL PRIMARY KEY,
      meaning_id INTEGER NOT NULL,
      text TEXT NOT NULL,
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
      FOREIGN KEY (meaning_id) REFERENCES meanings(id) ON DELETE CASCADE
  );
  ```
- Обновить модель и репозиторий

---

### Изменение 2: `words` query → `WordConnection`

**Текущая реализация:**
```go
func (r *queryResolver) Words(ctx context.Context, filter *WordFilter, limit *int, offset *int) ([]*Word, error)
```

**Новая реализация:**
```go
func (r *queryResolver) Words(ctx context.Context, filter *WordFilter, first *int, after *string) (*WordConnection, error) {
    // Парсинг cursor (base64 encoded offset)
    offset := 0
    if after != nil && *after != "" {
        // Декодировать cursor в offset
        // TODO: реализовать декодирование
    }
    
    // Нормализация first
    limit := 20
    if first != nil {
        limit = *first
    }
    
    // Получение данных
    words, err := r.words.List(ctx, ToWordFilter(filter), limit, offset)
    if err != nil {
        return nil, err
    }
    
    // Подсчет totalCount
    totalCount, err := r.words.Count(ctx, ToWordFilter(filter))
    if err != nil {
        return nil, err
    }
    
    // Создание edges
    edges := make([]*WordEdge, 0, len(words))
    for i, w := range words {
        cursor := encodeCursor(offset + i) // TODO: реализовать encodeCursor
        edges = append(edges, &WordEdge{
            Cursor: cursor,
            Node:   ToGraphQLWordBasic(&w),
        })
    }
    
    // Вычисление pageInfo
    hasNextPage := offset + len(words) < totalCount
    endCursor := ""
    if len(edges) > 0 {
        endCursor = edges[len(edges)-1].Cursor
    }
    
    return &WordConnection{
        Edges: edges,
        PageInfo: &PageInfo{
            HasNextPage: hasNextPage,
            EndCursor:   &endCursor,
        },
        TotalCount: totalCount,
    }, nil
}
```

**Нужно добавить в `converter.go`:**
```go
func encodeCursor(offset int) string {
    // Простая реализация: base64(offset)
    data := strconv.FormatInt(int64(offset), 10)
    return base64.StdEncoding.EncodeToString([]byte(data))
}

func decodeCursor(cursor string) (int, error) {
    data, err := base64.StdEncoding.DecodeString(cursor)
    if err != nil {
        return 0, err
    }
    offset, err := strconv.ParseInt(string(data), 10, 64)
    if err != nil {
        return 0, err
    }
    return int(offset), nil
}
```

---

### Изменение 3: Новые операции (Inbox, Suggest)

**Inbox операции требуют:**
1. Таблица `inbox_items`:
   ```sql
   CREATE TABLE inbox_items (
       id SERIAL PRIMARY KEY,
       text VARCHAR NOT NULL,
       source_context TEXT,
       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
   );
   ```

2. Модель:
   ```go
   type InboxItem struct {
       ID           int64     `db:"id"`
       Text         string    `db:"text"`
       SourceContext *string   `db:"source_context"`
       CreatedAt    time.Time `db:"created_at"`
   }
   ```

3. Резолверы (пример):
   ```go
   func (r *queryResolver) InboxItems(ctx context.Context) ([]*InboxItem, error) {
       items, err := r.inbox.List(ctx)
       if err != nil {
           return nil, transport.HandleError(ctx, err)
       }
       return ToGraphQLInboxItems(items), nil
   }
   
   func (r *mutationResolver) AddToInbox(ctx context.Context, text string, sourceContext *string) (*InboxItem, error) {
       item, err := r.inbox.Create(ctx, text, sourceContext)
       if err != nil {
           return nil, transport.HandleError(ctx, err)
       }
       return ToGraphQLInboxItem(item), nil
   }
   ```

**Suggest операция:**
- Требует интеграции с внешним API словаря (например, словарь.ру API)
- Или локальный поиск по существующим словам
- Возвращает `Suggestion` с `origin: LOCAL` или `DICTIONARY`

---

## ⚠️ Важные замечания

1. **Схема неполная:** Отсутствуют определения `CreateWordPayload` и `UpdateWordPayload` - их нужно добавить перед генерацией кода.

2. **Несовместимость БД:** `translationRu` как массив требует изменения структуры БД или временного решения.

3. **Frontend:** После изменений нужно обновить фронтенд:
   - Запросы `words` должны использовать `WordConnection`
   - `translationRu` теперь массив
   - Новые операции для Inbox

4. **Тесты:** Все существующие тесты нужно обновить под новую схему.

---

## 📚 Файлы, которые нужно изменить

### Обязательные изменения:
- `backend/graph/schema.graphqls` - добавить `CreateWordPayload`, `UpdateWordPayload`
- `backend/graph/schema.resolvers.go` - обновить все резолверы
- `backend/graph/converter.go` - обновить конвертеры
- `backend/graph/models_gen.go` - обновится после `make generate`
- `backend/graph/generated.go` - обновится после `make generate`

### Новые файлы (для Inbox):
- `backend/migrations/XXXXXX_create_inbox_items_table.sql`
- `backend/internal/model/inbox.go` (или добавить в `model.go`)
- `backend/internal/database/inbox/repo.go`
- `backend/internal/service/inbox/service.go`
- Обновить `backend/graph/resolver.go` для добавления `InboxService`

### Обновления существующих:
- `backend/internal/service/word/dto.go` - обновить `WordFilter`, `CreateWordInput`
- `backend/internal/service/word/list.go` - добавить фильтрацию по `status` и `tags`
- `backend/internal/database/word/` - возможно, добавить метод `Count()`

