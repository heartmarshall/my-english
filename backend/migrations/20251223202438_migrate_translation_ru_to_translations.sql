-- +goose Up
-- Миграция: делаем translation_ru в meanings nullable и опциональным
-- После миграции данных в translations, это поле можно будет удалить в будущем

-- Делаем поле nullable (если оно еще не nullable)
ALTER TABLE meanings ALTER COLUMN translation_ru DROP NOT NULL;

-- Комментарий для документации
COMMENT ON COLUMN meanings.translation_ru IS 'DEPRECATED: Используйте таблицу translations. Это поле оставлено для обратной совместимости и будет удалено в будущем.';

-- +goose Down
-- Возвращаем NOT NULL (если нужно)
-- ВАЖНО: Это может вызвать ошибки, если есть NULL значения

ALTER TABLE meanings ALTER COLUMN translation_ru SET NOT NULL;

