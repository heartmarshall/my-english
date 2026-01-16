package dataloader

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/heartmarshall/my-english/internal/database/repository"
	"github.com/heartmarshall/my-english/internal/model"
	"github.com/vikstrous/dataloadgen"
)

type ctxKey string

const loadersKey ctxKey = "dataloaders"

// Loaders содержит все доступные загрузчики данных.
type Loaders struct {
	// 1:N Loaders (Одно слово -> Много смыслов/картинок/произношений)
	SensesByEntryID         *dataloadgen.Loader[uuid.UUID, []model.Sense]
	ImagesByEntryID         *dataloadgen.Loader[uuid.UUID, []model.Image]
	PronunciationsByEntryID *dataloadgen.Loader[uuid.UUID, []model.Pronunciation]

	// 1:N Loaders (Один смысл -> Много примеров/переводов)
	ExamplesBySenseID     *dataloadgen.Loader[uuid.UUID, []model.Example]
	TranslationsBySenseID *dataloadgen.Loader[uuid.UUID, []model.Translation]

	// 1:1 Loaders (Одно слово -> Одна карточка)
	CardByEntryID *dataloadgen.Loader[uuid.UUID, *model.Card]
}

// NewLoaders создает экземпляры всех загрузчиков.
func NewLoaders(repos *repository.Registry) *Loaders {
	return &Loaders{
		SensesByEntryID: dataloadgen.NewLoader(
			newSensesByEntryIDFetcher(repos.Senses),
			dataloadgen.WithWait(time.Millisecond*2),
		),
		ImagesByEntryID: dataloadgen.NewLoader(
			newImagesByEntryIDFetcher(repos.Images),
			dataloadgen.WithWait(time.Millisecond*2),
		),
		PronunciationsByEntryID: dataloadgen.NewLoader(
			newPronunciationsByEntryIDFetcher(repos.Pronunciations),
			dataloadgen.WithWait(time.Millisecond*2),
		),
		ExamplesBySenseID: dataloadgen.NewLoader(
			newExamplesBySenseIDFetcher(repos.Examples),
			dataloadgen.WithWait(time.Millisecond*2),
		),
		TranslationsBySenseID: dataloadgen.NewLoader(
			newTranslationsBySenseIDFetcher(repos.Translations),
			dataloadgen.WithWait(time.Millisecond*2),
		),
		CardByEntryID: dataloadgen.NewLoader(
			newCardByEntryIDFetcher(repos.Cards),
			dataloadgen.WithWait(time.Millisecond*2),
		),
	}
}

// Middleware внедряет Loaders в контекст запроса.
// Loaders создаются заново для каждого запроса (Request-scoped),
// чтобы кэширование жило только в рамках одного запроса.
func Middleware(repos *repository.Registry) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			loaders := NewLoaders(repos)
			ctx := context.WithValue(r.Context(), loadersKey, loaders)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// For извлекает Loaders из контекста.
func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}
