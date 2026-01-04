# План реализации репозиториев

## Анализ структуры данных

### SYSTEM LAYER (Справочники)
- **DataSource** - источники данных (freedict, user, system)
  - Обычно read-only, заполняется при миграции
  - Приоритет: Низкий (можно использовать простой Base)

### LINGUISTIC LAYER (Глобальный словарь)

#### Основные сущности:
1. **Lexeme** (лексема) - основная сущность словаря
   - Приоритет: **ВЫСОКИЙ** ⭐
   - Методы: GetByID, FindByText, Search (триграммы), Create, Update, Delete
   - Связи: один-ко-многим с Pronunciations, Senses, Inflections

2. **Sense** (смысл/значение)
   - Приоритет: **ВЫСОКИЙ** ⭐
   - Методы: GetByID, ListByLexemeID, ListByPartOfSpeech, Create, Update, Delete
   - Связи: один-ко-многим с Translations, Examples, Relations

3. **Pronunciation** (произношение)
   - Приоритет: **СРЕДНИЙ**
   - Методы: GetByID, ListByLexemeID, GetByLexemeAndRegion, Create, Delete
   - Связи: многие-к-одному с Lexeme

4. **SenseTranslation** (перевод смысла)
   - Приоритет: **СРЕДНИЙ**
   - Методы: GetByID, ListBySenseID, Create, Delete
   - Связи: многие-к-одному с Sense

5. **Example** (пример использования)
   - Приоритет: **СРЕДНИЙ**
   - Методы: GetByID, ListBySenseID, Create, Update, Delete
   - Связи: многие-к-одному с Sense

6. **Inflection** (морфологическая связь)
   - Приоритет: **НИЗКИЙ**
   - Методы: GetByInflectedID, GetByLemmaID, Create, Delete
   - Связи: many-to-many между Lexemes (composite PK)

7. **SenseRelation** (семантическая связь)
   - Приоритет: **НИЗКИЙ**
   - Методы: GetBySourceSenseID, GetByTargetSenseID, GetByType, Create, Delete
   - Связи: many-to-many между Senses (composite PK)

### USER LAYER (Пользовательские данные)

#### Основные сущности:
1. **Card** (карточка пользователя)
   - Приоритет: **ВЫСОКИЙ** ⭐
   - Методы: GetByID, ListActive, ListBySenseID, ListDueForReview, Create, Update, SoftDelete
   - Связи: один-к-одному с SRSState, многие-ко-многим с Tags

2. **SRSState** (состояние SRS)
   - Приоритет: **ВЫСОКИЙ** ⭐
   - Методы: GetByCardID, ListDueForReview, ListByStatus, Create, Update, Delete
   - Связи: один-к-одному с Card (CardID = PK)

3. **ReviewLog** (лог повторений)
   - Приоритет: **ВЫСОКИЙ** ⭐
   - Методы: GetByID, ListByCardID, ListRecent, Create
   - Связи: многие-к-одному с Card

4. **Tag** (тег)
   - Приоритет: **СРЕДНИЙ**
   - Методы: GetByID, FindByName, ListAll, Create, Update, Delete
   - Связи: многие-ко-многим с Cards

5. **CardTag** (связь карточки и тега)
   - Приоритет: **СРЕДНИЙ**
   - Методы: GetByCardID, GetByTagID, Create, Delete
   - Связи: many-to-many (composite PK)

6. **InboxItem** (элемент inbox)
   - Приоритет: **СРЕДНИЙ**
   - Методы: GetByID, ListRecent, Create, Delete
   - Связи: нет

## План реализации (по приоритетам)

### Фаза 1: Критичные репозитории (ВЫСОКИЙ приоритет) ⭐

1. **LexemeRepository** ✅ (уже есть пример)
   - GetByID
   - FindByText (по text_normalized)
   - Search (триграммы, ILIKE)
   - Create
   - Update
   - Delete

2. **SenseRepository**
   - GetByID
   - ListByLexemeID
   - ListByPartOfSpeech
   - ListByCefrLevel
   - Create
   - Update
   - Delete

3. **CardRepository** ✅ (уже есть пример)
   - GetByID (только активные)
   - ListActive
   - ListBySenseID
   - ListDueForReview (JOIN с srs_states)
   - Create
   - Update
   - SoftDelete

4. **SRSStateRepository**
   - GetByCardID (CardID = PK)
   - ListDueForReview (WHERE due_date <= NOW())
   - ListByStatus
   - CreateOrUpdate (UPSERT)
   - Delete

5. **ReviewLogRepository**
   - GetByID
   - ListByCardID (сортировка по reviewed_at DESC)
   - ListRecent (последние N для всех карточек)
   - Create

### Фаза 2: Важные репозитории (СРЕДНИЙ приоритет)

6. **PronunciationRepository**
   - GetByID
   - ListByLexemeID
   - GetByLexemeAndRegion
   - Create
   - Delete

7. **SenseTranslationRepository**
   - GetByID
   - ListBySenseID
   - Create
   - Delete

8. **ExampleRepository**
   - GetByID
   - ListBySenseID
   - Create
   - Update
   - Delete

9. **TagRepository**
   - GetByID
   - FindByName
   - ListAll
   - Create
   - Update
   - Delete

10. **CardTagRepository**
    - GetByCardID (список тегов карточки)
    - GetByTagID (список карточек с тегом)
    - Create
    - Delete
    - DeleteByCardID (удалить все теги карточки)

11. **InboxRepository** ✅ (уже есть пример)
    - GetByID
    - ListRecent
    - Create
    - Delete

### Фаза 3: Дополнительные репозитории (НИЗКИЙ приоритет)

12. **InflectionRepository**
    - GetByInflectedID
    - GetByLemmaID
    - Create
    - Delete

13. **SenseRelationRepository**
    - GetBySourceSenseID
    - GetByTargetSenseID
    - GetByType
    - Create
    - Delete

14. **DataSourceRepository** (опционально)
    - GetByID
    - GetBySlug
    - ListAll

## Специальные методы для сложных запросов

### CardRepository
- `ListDueForReview()` - JOIN с srs_states, фильтр по due_date
- `ListWithSRSState()` - JOIN для получения карточек со статусом SRS
- `ListByTags()` - через CardTag, фильтр по нескольким тегам

### SenseRepository
- `ListByLexemeIDWithTranslations()` - JOIN с sense_translations
- `ListByLexemeIDWithExamples()` - JOIN с examples

### LexemeRepository
- `SearchFuzzy()` - триграммы для нечёткого поиска
- `GetWithPronunciations()` - JOIN с pronunciations
- `GetWithSenses()` - JOIN с senses

## Структура файлов

```
internal/database/repository/
├── doc.go                    # Документация
├── interfaces.go             # Интерфейсы
├── options.go                # Query options
├── repository.go             # Base[T]
├── example_repo.go          # Примеры (удалить после реализации)
│
├── lexeme_repo.go           # LexemeRepository
├── sense_repo.go            # SenseRepository
├── pronunciation_repo.go    # PronunciationRepository
├── sense_translation_repo.go # SenseTranslationRepository
├── example_repo.go          # ExampleRepository (переименовать)
├── inflection_repo.go       # InflectionRepository
├── sense_relation_repo.go   # SenseRelationRepository
│
├── card_repo.go             # CardRepository
├── srs_state_repo.go        # SRSStateRepository
├── review_log_repo.go       # ReviewLogRepository
├── tag_repo.go              # TagRepository
├── card_tag_repo.go         # CardTagRepository
├── inbox_repo.go            # InboxRepository
│
└── data_source_repo.go      # DataSourceRepository (опционально)
```

## Особенности реализации

### Composite Primary Keys
Для `Inflection`, `SenseRelation`, `CardTag`:
- Использовать `GetByCompositeKey()` вместо `GetByID()`
- `Delete()` принимает оба ключа

### JSONB поля
Для `SRSState.AlgorithmData`, `ReviewLog.StateBefore/StateAfter`:
- Использовать `pgx/v5` JSONB сканирование
- Возможно нужны helper методы для работы с JSONB

### Soft Delete
Для `Card`:
- Все методы чтения должны фильтровать `is_deleted = false`
- `Delete()` → `SoftDelete()` (установка флага)
- Возможно нужен `HardDelete()` для админки

### Триграммы поиск
Для `Lexeme`:
- Использовать `pg_trgm` расширение
- Метод `SearchFuzzy()` с `%similarity%` или `similarity()`

### UPSERT операции
Для `SRSState`:
- `CreateOrUpdate()` использует `ON CONFLICT (card_id) DO UPDATE`

## Тестирование

Каждый репозиторий должен иметь:
- Unit тесты с моками
- Integration тесты с testcontainers
- Тесты на транзакции
- Тесты на edge cases (NotFound, Duplicate, etc.)

