package dictionary

import (
	"context"
)

// Provider (бывший Fetcher) — интерфейс для внешнего источника словаря.
type Provider interface {
	// SourceSlug возвращает уникальный идентификатор источника (должен совпадать с data_sources.slug).
	// Например: "freedict", "cambridge", "gpt4".
	SourceSlug() string

	// Fetch получает данные о слове.
	// Возвращает nil, nil, если слово не найдено.
	Fetch(ctx context.Context, query string) (*ImportedWord, error)
}
