-- +goose Up
-- Включение расширения pg_trgm для триграммного поиска
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Создание GIN индекса на поле text для быстрого триграммного поиска
CREATE INDEX IF NOT EXISTS idx_words_text_trgm ON words USING gin (text gin_trgm_ops);

-- +goose Down
-- Удаление индекса и расширения
DROP INDEX IF EXISTS idx_words_text_trgm;
-- Расширение не удаляем, так как оно может использоваться другими таблицами

