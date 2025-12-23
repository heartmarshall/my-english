-- +goose Up
-- Создание таблицы связки Many-to-Many meanings_tags

CREATE TABLE meanings_tags (
    meaning_id INTEGER NOT NULL,
    tag_id INTEGER NOT NULL,
    
    PRIMARY KEY (meaning_id, tag_id),
    CONSTRAINT fk_meanings_tags_meaning_id FOREIGN KEY (meaning_id) REFERENCES meanings(id) ON DELETE CASCADE,
    CONSTRAINT fk_meanings_tags_tag_id FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

-- Индексы
CREATE INDEX idx_meanings_tags_meaning_id ON meanings_tags(meaning_id);
CREATE INDEX idx_meanings_tags_tag_id ON meanings_tags(tag_id);

-- +goose Down
-- Удаление таблицы meanings_tags

DROP TABLE IF EXISTS meanings_tags;

