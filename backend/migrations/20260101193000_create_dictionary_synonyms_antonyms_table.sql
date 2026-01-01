-- +goose Up
-- Создание таблицы dictionary_synonyms_antonyms для связывания значений словаря
-- Синонимы и антонимы относятся к конкретным значениям (dictionary_meanings), а не к словам целиком

CREATE TYPE relation_type AS ENUM ('synonym', 'antonym');

CREATE TABLE dictionary_synonyms_antonyms (
    id SERIAL PRIMARY KEY,
    meaning_id_1 INTEGER NOT NULL REFERENCES dictionary_meanings(id) ON DELETE CASCADE,
    meaning_id_2 INTEGER NOT NULL REFERENCES dictionary_meanings(id) ON DELETE CASCADE,
    relation_type relation_type NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Обеспечиваем уникальность связи (независимо от порядка meaning_id_1 и meaning_id_2)
    -- Используем CHECK для гарантии, что meaning_id_1 < meaning_id_2
    CONSTRAINT check_meaning_order CHECK (meaning_id_1 < meaning_id_2),
    CONSTRAINT unique_relation UNIQUE (meaning_id_1, meaning_id_2, relation_type)
);

-- Индексы для быстрого поиска
CREATE INDEX idx_dictionary_synonyms_antonyms_meaning_1 ON dictionary_synonyms_antonyms(meaning_id_1);
CREATE INDEX idx_dictionary_synonyms_antonyms_meaning_2 ON dictionary_synonyms_antonyms(meaning_id_2);
CREATE INDEX idx_dictionary_synonyms_antonyms_type ON dictionary_synonyms_antonyms(relation_type);
-- Составной индекс для поиска всех связей для значения
CREATE INDEX idx_dictionary_synonyms_antonyms_meaning_1_type ON dictionary_synonyms_antonyms(meaning_id_1, relation_type);
CREATE INDEX idx_dictionary_synonyms_antonyms_meaning_2_type ON dictionary_synonyms_antonyms(meaning_id_2, relation_type);

-- Триггер для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_dictionary_synonyms_antonyms_updated_at() RETURNS TRIGGER LANGUAGE plpgsql AS 'BEGIN NEW.updated_at = CURRENT_TIMESTAMP; RETURN NEW; END;';

CREATE TRIGGER update_dictionary_synonyms_antonyms_updated_at BEFORE UPDATE ON dictionary_synonyms_antonyms
    FOR EACH ROW EXECUTE FUNCTION update_dictionary_synonyms_antonyms_updated_at();

-- +goose Down
-- Удаление таблицы dictionary_synonyms_antonyms

DROP TRIGGER IF EXISTS update_dictionary_synonyms_antonyms_updated_at ON dictionary_synonyms_antonyms;
DROP FUNCTION IF EXISTS update_dictionary_synonyms_antonyms_updated_at();
DROP INDEX IF EXISTS idx_dictionary_synonyms_antonyms_meaning_2_type;
DROP INDEX IF EXISTS idx_dictionary_synonyms_antonyms_meaning_1_type;
DROP INDEX IF EXISTS idx_dictionary_synonyms_antonyms_type;
DROP INDEX IF EXISTS idx_dictionary_synonyms_antonyms_meaning_2;
DROP INDEX IF EXISTS idx_dictionary_synonyms_antonyms_meaning_1;
DROP TABLE IF EXISTS dictionary_synonyms_antonyms;
DROP TYPE IF EXISTS relation_type;

