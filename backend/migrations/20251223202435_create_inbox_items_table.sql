-- +goose Up
-- Создание таблицы inbox_items

CREATE TABLE inbox_items (
    id SERIAL PRIMARY KEY,
    text VARCHAR NOT NULL,
    source_context TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Индекс для сортировки по дате создания
CREATE INDEX idx_inbox_items_created_at ON inbox_items(created_at);

-- Индекс для поиска по тексту
CREATE INDEX idx_inbox_items_text ON inbox_items(text);

-- +goose Down
-- Удаление таблицы inbox_items

DROP INDEX IF EXISTS idx_inbox_items_text;
DROP INDEX IF EXISTS idx_inbox_items_created_at;
DROP TABLE IF EXISTS inbox_items;

