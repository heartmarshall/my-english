-- +goose Up
-- Создание таблицы examples

CREATE TABLE examples (
    id SERIAL PRIMARY KEY,
    meaning_id INTEGER NOT NULL,
    sentence_en TEXT NOT NULL,
    sentence_ru TEXT,
    source_name example_source,
    
    CONSTRAINT fk_examples_meaning_id FOREIGN KEY (meaning_id) REFERENCES meanings(id) ON DELETE CASCADE
);

-- Индексы
CREATE INDEX idx_examples_meaning_id ON examples(meaning_id);
CREATE INDEX idx_examples_source_name ON examples(source_name);

-- +goose Down
-- Удаление таблицы examples

DROP TABLE IF EXISTS examples;

