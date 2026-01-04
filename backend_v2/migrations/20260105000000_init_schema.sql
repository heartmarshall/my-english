-- +goose Up
-- ============================================================================
-- INIT V2 SCHEMA
-- Полная перезагрузка схемы данных под архитектуру Dictionary + User Progress
-- ============================================================================

-- 1. Включаем необходимые расширения
CREATE EXTENSION IF NOT EXISTS "pgcrypto"; -- Для gen_random_uuid()
CREATE EXTENSION IF NOT EXISTS "pg_trgm";  -- Для нечеткого поиска

-- 2. Очистка старой схемы (на случай, если она была)
-- Удаляем таблицы в обратном порядке зависимости
DROP TABLE IF EXISTS card_tags CASCADE;
DROP TABLE IF EXISTS review_logs CASCADE;
DROP TABLE IF EXISTS srs_states CASCADE;
DROP TABLE IF EXISTS cards CASCADE;
DROP TABLE IF EXISTS inbox_items CASCADE;
DROP TABLE IF EXISTS examples CASCADE;
DROP TABLE IF EXISTS sense_relations CASCADE;
DROP TABLE IF EXISTS sense_translations CASCADE;
DROP TABLE IF EXISTS senses CASCADE;
DROP TABLE IF EXISTS inflections CASCADE;
DROP TABLE IF EXISTS pronunciations CASCADE;
DROP TABLE IF EXISTS lexemes CASCADE;
DROP TABLE IF EXISTS data_sources CASCADE;
DROP TABLE IF EXISTS tags CASCADE;

-- Удаляем старые типы, если они были
DROP TYPE IF EXISTS morphological_type CASCADE;
DROP TYPE IF EXISTS relation_type CASCADE;
DROP TYPE IF EXISTS accent_region CASCADE;
DROP TYPE IF EXISTS part_of_speech CASCADE;
DROP TYPE IF EXISTS learning_status CASCADE;

-- 3. Создание ENUM типов
CREATE TYPE learning_status AS ENUM ('new', 'learning', 'review', 'mastered');
CREATE TYPE part_of_speech AS ENUM ('noun', 'verb', 'adjective', 'adverb', 'pronoun', 'preposition', 'conjunction', 'interjection', 'phrase', 'idiom', 'other');
CREATE TYPE accent_region AS ENUM ('us', 'uk', 'au', 'general');
CREATE TYPE relation_type AS ENUM ('synonym', 'antonym', 'related', 'collocation');
CREATE TYPE morphological_type AS ENUM ('plural', 'past_tense', 'past_participle', 'present_participle', 'comparative', 'superlative');

-- 4. SYSTEM LAYER (Источники данных)
CREATE TABLE data_sources (
    id SERIAL PRIMARY KEY,
    slug VARCHAR(50) NOT NULL UNIQUE, -- 'freedict', 'user', 'gpt4'
    display_name VARCHAR(100) NOT NULL,
    trust_level INTEGER DEFAULT 5,    -- 1-10, где 10 - максимальное доверие (юзер)
    website_url VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Заполним дефолтные источники
INSERT INTO data_sources (slug, display_name, trust_level) VALUES 
('user', 'User Manual Entry', 10),
('freedict', 'Free Dictionary API', 8),
('system', 'System Import', 9);

-- 5. LINGUISTIC LAYER (Глобальный словарь)

-- Лексемы (Слова)
CREATE TABLE lexemes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Нормализованный текст для поиска (lowercase, trim)
    text_normalized VARCHAR(100) NOT NULL UNIQUE,
    -- Оригинальное отображение (например "London" с большой буквы)
    text_display VARCHAR(100) NOT NULL,
    
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_lexemes_text_trgm ON lexemes USING gin (text_normalized gin_trgm_ops);

-- Произношение
CREATE TABLE pronunciations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lexeme_id UUID NOT NULL REFERENCES lexemes(id) ON DELETE CASCADE,
    
    audio_url VARCHAR(2048) NOT NULL,
    transcription VARCHAR(100),
    region accent_region DEFAULT 'general',
    
    source_id INTEGER REFERENCES data_sources(id)
);

CREATE INDEX idx_pronunciations_lexeme ON pronunciations(lexeme_id);

-- Морфология (Связи форм слов)
CREATE TABLE inflections (
    inflected_lexeme_id UUID REFERENCES lexemes(id) ON DELETE CASCADE,
    lemma_lexeme_id UUID REFERENCES lexemes(id) ON DELETE CASCADE,
    type morphological_type NOT NULL,
    
    PRIMARY KEY (inflected_lexeme_id, lemma_lexeme_id)
);

-- Смыслы (Meanings/Senses)
CREATE TABLE senses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lexeme_id UUID NOT NULL REFERENCES lexemes(id) ON DELETE CASCADE,
    
    part_of_speech part_of_speech NOT NULL,
    definition TEXT NOT NULL,
    cefr_level VARCHAR(2), -- A1, B2...
    
    source_id INTEGER NOT NULL REFERENCES data_sources(id),
    external_ref_id VARCHAR(100), -- ID во внешней системе для обновлений
    
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_senses_lexeme ON senses(lexeme_id);

-- Переводы смыслов (Словарь может давать несколько вариантов)
CREATE TABLE sense_translations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sense_id UUID NOT NULL REFERENCES senses(id) ON DELETE CASCADE,
    translation TEXT NOT NULL,
    source_id INTEGER REFERENCES data_sources(id)
);

CREATE INDEX idx_sense_translations_sense ON sense_translations(sense_id);

-- Семантические связи (Синонимы и т.д.)
CREATE TABLE sense_relations (
    source_sense_id UUID NOT NULL REFERENCES senses(id) ON DELETE CASCADE,
    target_sense_id UUID NOT NULL REFERENCES senses(id) ON DELETE CASCADE,
    type relation_type NOT NULL,
    
    -- Направленность связи (true = двусторонняя, false = направленная source->target)
    is_bidirectional BOOLEAN NOT NULL DEFAULT TRUE,
    
    source_id INTEGER REFERENCES data_sources(id),
    
    PRIMARY KEY (source_sense_id, target_sense_id, type)
);

-- Примеры предложений
CREATE TABLE examples (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sense_id UUID REFERENCES senses(id) ON DELETE CASCADE,
    
    sentence_en TEXT NOT NULL,
    sentence_ru TEXT,
    
    -- Индексы [start, end] для подсветки целевого слова в sentence_en
    -- Пример: "I went home", went=[2, 6]
    target_word_range INTEGER[], 
    
    source_name VARCHAR(255) -- "Harry Potter, ch. 4"
);

CREATE INDEX idx_examples_sense ON examples(sense_id);

-- 6. USER LAYER (Пользовательские данные)

-- Inbox (GTD)
CREATE TABLE inbox_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    raw_text VARCHAR(255) NOT NULL,
    context_note TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Теги (Категории)
CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    color_hex VARCHAR(7)
);

-- Личные карточки (Cards)
CREATE TABLE cards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Ссылка на словарь (Может быть NULL, если слово полностью кастомное)
    sense_id UUID REFERENCES senses(id) ON DELETE SET NULL,
    
    -- User Overrides (Кастомные данные пользователя)
    custom_text VARCHAR(100),
    custom_transcription VARCHAR(100),
    custom_translations TEXT[], -- Массив строк для удобства
    custom_note TEXT,           -- Личные заметки (Markdown)
    custom_image_url VARCHAR(2048),
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    is_deleted BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_cards_sense ON cards(sense_id);
CREATE INDEX idx_cards_created_at ON cards(created_at);

-- Связь Карточки и Теги
CREATE TABLE card_tags (
    card_id UUID REFERENCES cards(id) ON DELETE CASCADE,
    tag_id INTEGER REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (card_id, tag_id)
);

-- SRS State (Текущий прогресс)
CREATE TABLE srs_states (
    card_id UUID PRIMARY KEY REFERENCES cards(id) ON DELETE CASCADE,
    
    status learning_status NOT NULL DEFAULT 'new',
    due_date TIMESTAMPTZ, -- Когда повторять
    
    -- Данные алгоритма (FSRS / SM-2)
    -- JSONB позволяет хранить stability, difficulty и т.д. без изменения схемы
    algorithm_data JSONB NOT NULL DEFAULT '{}'::jsonb,
    
    last_review_at TIMESTAMPTZ
);

CREATE INDEX idx_srs_states_due_date ON srs_states(due_date);
CREATE INDEX idx_srs_states_status ON srs_states(status);

-- Review Logs (История для аналитики и обучения алгоритма)
CREATE TABLE review_logs (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    card_id UUID NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
    
    grade SMALLINT NOT NULL CHECK (grade >= 1 AND grade <= 5),
    duration_ms INTEGER, -- Время ответа в миллисекундах
    
    reviewed_at TIMESTAMPTZ DEFAULT NOW(),
    
    -- Снапшоты состояния алгоритма (для отладки)
    state_before JSONB,
    state_after JSONB
);

CREATE INDEX idx_review_logs_card_id ON review_logs(card_id);
CREATE INDEX idx_review_logs_reviewed_at ON review_logs(reviewed_at);


-- +goose Down
-- ============================================================================
-- ROLLBACK
-- ============================================================================

DROP TABLE IF EXISTS review_logs;
DROP TABLE IF EXISTS srs_states;
DROP TABLE IF EXISTS card_tags;
DROP TABLE IF EXISTS cards;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS inbox_items;
DROP TABLE IF EXISTS examples;
DROP TABLE IF EXISTS sense_relations;
DROP TABLE IF EXISTS sense_translations;
DROP TABLE IF EXISTS senses;
DROP TABLE IF EXISTS inflections;
DROP TABLE IF EXISTS pronunciations;
DROP TABLE IF EXISTS lexemes;
DROP TABLE IF EXISTS data_sources;

DROP TYPE IF EXISTS morphological_type;
DROP TYPE IF EXISTS relation_type;
DROP TYPE IF EXISTS accent_region;
DROP TYPE IF EXISTS part_of_speech;
DROP TYPE IF EXISTS learning_status;