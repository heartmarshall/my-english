-- +goose Up
-- Создание таблицы translations для хранения множественных переводов значений

CREATE TABLE translations (
    id SERIAL PRIMARY KEY,
    meaning_id INTEGER NOT NULL REFERENCES meanings(id) ON DELETE CASCADE,
    translation_ru TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Уникальность: одно значение не может иметь два одинаковых перевода
    UNIQUE(meaning_id, translation_ru)
);

-- Индекс для быстрого поиска переводов по meaning_id
CREATE INDEX idx_translations_meaning_id ON translations(meaning_id);

-- Индекс для поиска по тексту перевода (может быть полезен для поиска)
CREATE INDEX idx_translations_translation_ru ON translations(translation_ru);

-- Миграция данных: переносим существующие translation_ru из meanings в translations
-- Создаем записи для всех существующих meanings
INSERT INTO translations (meaning_id, translation_ru, created_at)
SELECT id, translation_ru, created_at
FROM meanings
WHERE translation_ru IS NOT NULL AND translation_ru != '';

-- +goose Down
-- Удаление таблицы translations
-- ВАЖНО: Данные из meanings.translation_ru уже потеряны, так как мы их перенесли
-- Если нужно восстановить, можно добавить обратную миграцию данных

DROP INDEX IF EXISTS idx_translations_translation_ru;
DROP INDEX IF EXISTS idx_translations_meaning_id;
DROP TABLE IF EXISTS translations;

