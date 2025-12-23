# Миграции базы данных

Этот каталог содержит миграции базы данных, управляемые с помощью [goose](https://github.com/pressly/goose).

## Структура миграций

Миграции следуют формату `YYYYMMDDHHMMSS_description.sql` и используют директивы goose:
- `-- +goose Up` - код для применения миграции
- `-- +goose Down` - код для отката миграции

Оба блока находятся в одном файле.

## Использование

### Применить все миграции
```bash
make migrate-up DB_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable
```

### Откатить последнюю миграцию
```bash
make migrate-down DB_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable
```

### Проверить статус миграций
```bash
make migrate-status DB_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable
```

### Создать новую миграцию
```bash
make migrate-create NAME=add_new_feature
```

## Порядок миграций

1. `20251223202429_create_enums` - создание ENUM типов
2. `20251223202430_create_words_table` - создание таблицы words
3. `20251223202431_create_meanings_table` - создание таблицы meanings
4. `20251223202432_create_examples_table` - создание таблицы examples
5. `20251223202433_create_tags_table` - создание таблицы tags
6. `20251223202434_create_meanings_tags_table` - создание таблицы связки meanings_tags

## Примечания

- Все миграции должны быть обратимыми (иметь соответствующий `.down.sql` файл)
- При создании новых миграций используйте команду `make migrate-create`
- Перед применением миграций убедитесь, что база данных существует

