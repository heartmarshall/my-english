-- +goose Up
-- Создание ENUM типов

-- Статус изучения слова
CREATE TYPE learning_status AS ENUM ('new', 'learning', 'review', 'mastered');

-- Часть речи
CREATE TYPE part_of_speech AS ENUM ('noun', 'verb', 'adjective', 'adverb', 'other');

-- Источник примера
CREATE TYPE example_source AS ENUM ('film', 'book', 'chat', 'video', 'podcast');

-- +goose Down
-- Удаление ENUM типов

DROP TYPE IF EXISTS example_source;
DROP TYPE IF EXISTS part_of_speech;
DROP TYPE IF EXISTS learning_status;

