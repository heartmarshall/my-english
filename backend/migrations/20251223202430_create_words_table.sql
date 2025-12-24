-- +goose Up
-- Создание таблицы words

CREATE TABLE words (
    id SERIAL PRIMARY KEY,
    text VARCHAR NOT NULL UNIQUE,
    transcription VARCHAR,
    audio_url VARCHAR,
    frequency_rank INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Индекс для быстрого поиска по тексту
CREATE INDEX idx_words_text ON words(text);

-- Индекс для сортировки по дате создания
CREATE INDEX idx_words_created_at ON words(created_at);

-- +goose Down
DROP TABLE IF EXISTS words;

