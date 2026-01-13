package suggestion

import (
	"context"

	"github.com/heartmarshall/my-english/internal/service/dictionary"
)

// Result — унифицированный ответ от провайдера.
type Result struct {
	SourceSlug     string
	SourceName     string
	Senses         []dictionary.SenseInput // Используем те же Input структуры, что и для создания слова
	Images         []dictionary.ImageInput
	Pronunciations []dictionary.PronunciationInput
}

// Provider — интерфейс внешнего источника данных.
type Provider interface {
	// Slug возвращает уникальный идентификатор провайдера ("freedict", "openai").
	Slug() string
	// Name возвращает человекочитаемое название.
	Name() string
	// Fetch запрашивает данные о слове.
	Fetch(ctx context.Context, text string) (*Result, error)
}
