-- +goose Up
-- Создание таблицы dictionary_meanings для значений слов из внутреннего словаря

CREATE TABLE dictionary_meanings (
    id SERIAL PRIMARY KEY,
    dictionary_word_id INTEGER NOT NULL REFERENCES dictionary_words(id) ON DELETE CASCADE,
    
    -- Лингвистика
    part_of_speech part_of_speech NOT NULL,
    definition_en TEXT,
    cefr_level VARCHAR,
    image_url VARCHAR,
    
    -- Порядок значения (для сортировки)
    order_index INTEGER NOT NULL DEFAULT 0,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_dictionary_meanings_word_id FOREIGN KEY (dictionary_word_id) REFERENCES dictionary_words(id) ON DELETE CASCADE
);

-- Индексы
CREATE INDEX idx_dictionary_meanings_word_id ON dictionary_meanings(dictionary_word_id);
CREATE INDEX idx_dictionary_meanings_part_of_speech ON dictionary_meanings(part_of_speech);
CREATE INDEX idx_dictionary_meanings_order_index ON dictionary_meanings(dictionary_word_id, order_index);

-- Триггер для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_dictionary_meanings_updated_at() RETURNS TRIGGER LANGUAGE plpgsql AS 'BEGIN NEW.updated_at = CURRENT_TIMESTAMP; RETURN NEW; END;';

CREATE TRIGGER update_dictionary_meanings_updated_at BEFORE UPDATE ON dictionary_meanings
    FOR EACH ROW EXECUTE FUNCTION update_dictionary_meanings_updated_at();

-- +goose Down
-- Удаление таблицы dictionary_meanings

DROP TRIGGER IF EXISTS update_dictionary_meanings_updated_at ON dictionary_meanings;
DROP FUNCTION IF EXISTS update_dictionary_meanings_updated_at();
DROP INDEX IF EXISTS idx_dictionary_meanings_order_index;
DROP INDEX IF EXISTS idx_dictionary_meanings_part_of_speech;
DROP INDEX IF EXISTS idx_dictionary_meanings_word_id;
DROP TABLE IF EXISTS dictionary_meanings;

