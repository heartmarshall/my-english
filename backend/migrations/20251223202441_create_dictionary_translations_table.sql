-- +goose Up
-- Создание таблицы dictionary_translations для переводов значений из внутреннего словаря

CREATE TABLE dictionary_translations (
    id SERIAL PRIMARY KEY,
    dictionary_meaning_id INTEGER NOT NULL REFERENCES dictionary_meanings(id) ON DELETE CASCADE,
    translation_ru TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Уникальность: одно значение не может иметь два одинаковых перевода
    UNIQUE(dictionary_meaning_id, translation_ru)
);

-- Индексы
CREATE INDEX idx_dictionary_translations_meaning_id ON dictionary_translations(dictionary_meaning_id);
CREATE INDEX idx_dictionary_translations_translation_ru ON dictionary_translations(translation_ru);

-- +goose Down
-- Удаление таблицы dictionary_translations

DROP INDEX IF EXISTS idx_dictionary_translations_translation_ru;
DROP INDEX IF EXISTS idx_dictionary_translations_meaning_id;
DROP TABLE IF EXISTS dictionary_translations;

