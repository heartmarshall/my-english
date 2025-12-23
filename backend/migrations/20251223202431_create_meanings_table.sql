-- +goose Up
-- Создание таблицы meanings

CREATE TABLE meanings (
    id SERIAL PRIMARY KEY,
    word_id INTEGER NOT NULL,
    
    -- Лингвистика
    part_of_speech part_of_speech NOT NULL,
    definition_en TEXT,
    translation_ru TEXT NOT NULL,
    cefr_level VARCHAR,
    image_url VARCHAR,
    
    -- SRS (Интервальное повторение)
    learning_status learning_status NOT NULL DEFAULT 'new',
    next_review_at TIMESTAMP,
    interval INTEGER,
    ease_factor DECIMAL(5, 2) DEFAULT 2.5,
    review_count INTEGER DEFAULT 0,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_meanings_word_id FOREIGN KEY (word_id) REFERENCES words(id) ON DELETE CASCADE
);

-- Индексы
CREATE INDEX idx_meanings_word_id ON meanings(word_id);
CREATE INDEX idx_meanings_learning_status ON meanings(learning_status);
CREATE INDEX idx_meanings_next_review_at ON meanings(next_review_at);
CREATE INDEX idx_meanings_part_of_speech ON meanings(part_of_speech);

-- Триггер для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_meanings_updated_at BEFORE UPDATE ON meanings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- +goose Down
-- Удаление таблицы meanings

DROP TRIGGER IF EXISTS update_meanings_updated_at ON meanings;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP TABLE IF EXISTS meanings;

