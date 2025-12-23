-- +goose Up
-- Создание таблицы tags

CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL UNIQUE
);

-- Индекс для быстрого поиска по имени тега
CREATE INDEX idx_tags_name ON tags(name);

-- +goose Down
-- Удаление таблицы tags

DROP TABLE IF EXISTS tags;

