-- +goose Up
-- Удаляем ограничение уникальности только по тексту
ALTER TABLE dictionary_words DROP CONSTRAINT IF EXISTS dictionary_words_text_key;
DROP INDEX IF EXISTS idx_dictionary_words_text;

-- Добавляем ограничение уникальности по паре (текст, источник)
CREATE UNIQUE INDEX idx_dictionary_words_text_source ON dictionary_words(text, source);

-- +goose Down
DROP INDEX IF EXISTS idx_dictionary_words_text_source;
CREATE UNIQUE INDEX idx_dictionary_words_text ON dictionary_words(text);
ALTER TABLE dictionary_words ADD CONSTRAINT dictionary_words_text_key UNIQUE (text);