-- +goose Up
-- Создание таблицы dictionary_word_forms для хранения различных форм слов
-- Формы слов: времена глаголов (go, went, gone), множественное число существительных (mouse, mice),
-- степени сравнения прилагательных (big, bigger, biggest) и т.д.

CREATE TABLE dictionary_word_forms (
    id SERIAL PRIMARY KEY,
    dictionary_word_id INTEGER NOT NULL REFERENCES dictionary_words(id) ON DELETE CASCADE,
    form_text VARCHAR NOT NULL,
    form_type VARCHAR, -- Тип формы: 'past_tense', 'past_participle', 'plural', 'comparative', 'superlative', 'third_person_singular', 'present_participle', 'gerund' и т.д.
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(dictionary_word_id, form_text, form_type)
);

-- Индексы для быстрого поиска
CREATE INDEX idx_dictionary_word_forms_dictionary_word_id ON dictionary_word_forms(dictionary_word_id);
CREATE INDEX idx_dictionary_word_forms_form_text ON dictionary_word_forms(form_text);
CREATE INDEX idx_dictionary_word_forms_form_type ON dictionary_word_forms(form_type);

-- GIN индекс для триграммного поиска по формам
CREATE INDEX idx_dictionary_word_forms_form_text_trgm ON dictionary_word_forms USING gin (form_text gin_trgm_ops);

-- Триггер для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_dictionary_word_forms_updated_at() RETURNS TRIGGER LANGUAGE plpgsql AS 'BEGIN NEW.updated_at = CURRENT_TIMESTAMP; RETURN NEW; END;';

CREATE TRIGGER update_dictionary_word_forms_updated_at BEFORE UPDATE ON dictionary_word_forms
    FOR EACH ROW EXECUTE FUNCTION update_dictionary_word_forms_updated_at();

-- +goose Down
-- Удаление таблицы dictionary_word_forms

DROP TRIGGER IF EXISTS update_dictionary_word_forms_updated_at ON dictionary_word_forms;
DROP FUNCTION IF EXISTS update_dictionary_word_forms_updated_at();
DROP INDEX IF EXISTS idx_dictionary_word_forms_form_text_trgm;
DROP INDEX IF EXISTS idx_dictionary_word_forms_form_type;
DROP INDEX IF EXISTS idx_dictionary_word_forms_form_text;
DROP INDEX IF EXISTS idx_dictionary_word_forms_dictionary_word_id;
DROP TABLE IF EXISTS dictionary_word_forms;

