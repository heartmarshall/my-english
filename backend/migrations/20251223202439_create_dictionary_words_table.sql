-- +goose Up
-- Создание таблицы dictionary_words для внутреннего словаря
-- Этот словарь используется для хранения данных из внешних источников (Free Dictionary API, Oxford и т.д.)
-- и не является частью пользовательского словаря

CREATE TABLE dictionary_words (
    id SERIAL PRIMARY KEY,
    text VARCHAR NOT NULL UNIQUE,
    transcription VARCHAR,
    audio_url VARCHAR,
    frequency_rank INTEGER,
    source VARCHAR NOT NULL, -- Источник данных: 'free_dictionary', 'oxford', 'custom' и т.д.
    source_id VARCHAR, -- ID слова в источнике (если есть)
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для быстрого поиска
CREATE INDEX idx_dictionary_words_text ON dictionary_words(text);
CREATE INDEX idx_dictionary_words_source ON dictionary_words(source);
CREATE INDEX idx_dictionary_words_created_at ON dictionary_words(created_at);

-- GIN индекс для триграммного поиска (как в words)
CREATE INDEX idx_dictionary_words_text_trgm ON dictionary_words USING gin (text gin_trgm_ops);

-- Триггер для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_dictionary_words_updated_at() RETURNS TRIGGER LANGUAGE plpgsql AS 'BEGIN NEW.updated_at = CURRENT_TIMESTAMP; RETURN NEW; END;';

CREATE TRIGGER update_dictionary_words_updated_at BEFORE UPDATE ON dictionary_words
    FOR EACH ROW EXECUTE FUNCTION update_dictionary_words_updated_at();

-- +goose Down
-- Удаление таблицы dictionary_words

DROP TRIGGER IF EXISTS update_dictionary_words_updated_at ON dictionary_words;
DROP FUNCTION IF EXISTS update_dictionary_words_updated_at();
DROP INDEX IF EXISTS idx_dictionary_words_text_trgm;
DROP INDEX IF EXISTS idx_dictionary_words_created_at;
DROP INDEX IF EXISTS idx_dictionary_words_source;
DROP INDEX IF EXISTS idx_dictionary_words_text;
DROP TABLE IF EXISTS dictionary_words;

