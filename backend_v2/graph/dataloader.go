package graph

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/heartmarshall/my-english/internal/model"
	"github.com/heartmarshall/my-english/internal/service/loader"
	"github.com/vikstrous/dataloadgen"
)

// Loaders содержит все доступные DataLoaders.
type Loaders struct {
	SensesByLexemeID         *dataloadgen.Loader[uuid.UUID, []model.Sense]
	PronunciationsByLexemeID *dataloadgen.Loader[uuid.UUID, []model.Pronunciation]
	TranslationsBySenseID    *dataloadgen.Loader[uuid.UUID, []model.SenseTranslation]
	ExamplesBySenseID        *dataloadgen.Loader[uuid.UUID, []model.Example]
	SRSByCardID              *dataloadgen.Loader[uuid.UUID, *model.SRSState] // 1:1 связь, поэтому указатель
	TagsByCardID             *dataloadgen.Loader[uuid.UUID, []model.Tag]
	SensesByID               *dataloadgen.Loader[uuid.UUID, *model.Sense]          // Для Card.sense
	LexemesByID              *dataloadgen.Loader[uuid.UUID, *model.Lexeme]         // Для Sense.lexeme
	DataSourcesByID          *dataloadgen.Loader[int, *model.DataSource]           // Для Source полей
	ReviewCountByCardID      *dataloadgen.Loader[uuid.UUID, int]                   // Для CardProgress.reviewCount
	RelationsBySenseID       *dataloadgen.Loader[uuid.UUID, []model.SenseRelation] // Для Sense.relations
}

// NewLoaders инициализирует лоадеры.
func NewLoaders(svc *loader.Service) *Loaders {
	return &Loaders{
		SensesByLexemeID: dataloadgen.NewLoader(
			buildGroupedBatcher(svc.GetSensesByLexemeIDs, func(item model.Sense) uuid.UUID { return item.LexemeID }),
			dataloadgen.WithWait(2*time.Millisecond),
		),
		PronunciationsByLexemeID: dataloadgen.NewLoader(
			buildGroupedBatcher(svc.GetPronunciationsByLexemeIDs, func(item model.Pronunciation) uuid.UUID { return item.LexemeID }),
			dataloadgen.WithWait(2*time.Millisecond),
		),
		TranslationsBySenseID: dataloadgen.NewLoader(
			buildGroupedBatcher(svc.GetTranslationsBySenseIDs, func(item model.SenseTranslation) uuid.UUID { return item.SenseID }),
			dataloadgen.WithWait(2*time.Millisecond),
		),
		ExamplesBySenseID: dataloadgen.NewLoader(
			buildGroupedBatcher(svc.GetExamplesBySenseIDs, func(item model.Example) uuid.UUID {
				if item.SenseID == nil {
					return uuid.Nil
				}
				return *item.SenseID
			}),
			dataloadgen.WithWait(2*time.Millisecond),
		),
		SRSByCardID: dataloadgen.NewLoader(
			buildOneToOneBatcher(svc.GetSRSByCardIDs, func(item model.SRSState) uuid.UUID { return item.CardID }),
			dataloadgen.WithWait(2*time.Millisecond),
		),
		TagsByCardID: dataloadgen.NewLoader(
			func(ctx context.Context, keys []uuid.UUID) ([][]model.Tag, []error) {
				// У тегов сложная логика в сервисе (возвращает Map), поэтому отдельная функция-адаптер
				tagsMap, err := svc.GetTagsByCardIDs(ctx, keys)
				if err != nil {
					// Возвращаем ошибку для всех ключей
					errs := make([]error, len(keys))
					for i := range errs {
						errs[i] = err
					}
					return nil, errs
				}

				// Мапим результаты в порядке ключей
				result := make([][]model.Tag, len(keys))
				for i, key := range keys {
					result[i] = tagsMap[key] // вернет nil (empty slice) если нет в мапе, это ок
				}
				return result, nil
			},
			dataloadgen.WithWait(2*time.Millisecond),
		),
		SensesByID: dataloadgen.NewLoader(
			buildOneToOneBatcher(svc.GetSensesByIDs, func(item model.Sense) uuid.UUID { return item.ID }),
			dataloadgen.WithWait(2*time.Millisecond),
		),
		LexemesByID: dataloadgen.NewLoader(
			buildOneToOneBatcher(svc.GetLexemesByIDs, func(item model.Lexeme) uuid.UUID { return item.ID }),
			dataloadgen.WithWait(2*time.Millisecond),
		),
		DataSourcesByID: dataloadgen.NewLoader(
			func(ctx context.Context, keys []int) ([]*model.DataSource, []error) {
				sources, err := svc.GetDataSourcesByIDs(ctx, keys)
				if err != nil {
					errs := make([]error, len(keys))
					for i := range errs {
						errs[i] = err
					}
					return nil, errs
				}

				// Мапим по ID
				mapped := make(map[int]*model.DataSource)
				for i := range sources {
					mapped[sources[i].ID] = &sources[i]
				}

				// Расставляем
				result := make([]*model.DataSource, len(keys))
				for i, key := range keys {
					if source, ok := mapped[key]; ok {
						result[i] = source
					} else {
						result[i] = nil
					}
				}
				return result, nil
			},
			dataloadgen.WithWait(2*time.Millisecond),
		),
		ReviewCountByCardID: dataloadgen.NewLoader(
			func(ctx context.Context, keys []uuid.UUID) ([]int, []error) {
				countsMap, err := svc.GetReviewCountByCardIDs(ctx, keys)
				if err != nil {
					errs := make([]error, len(keys))
					for i := range errs {
						errs[i] = err
					}
					return nil, errs
				}

				result := make([]int, len(keys))
				for i, key := range keys {
					result[i] = countsMap[key]
				}
				return result, nil
			},
			dataloadgen.WithWait(2*time.Millisecond),
		),
		RelationsBySenseID: dataloadgen.NewLoader(
			buildGroupedBatcher(svc.GetRelationsBySenseIDs, func(item model.SenseRelation) uuid.UUID { return item.SourceSenseID }),
			dataloadgen.WithWait(2*time.Millisecond),
		),
	}
}

// --- Helpers для создания batch-функций ---

// buildGroupedBatcher создает функцию для One-to-Many связей (например, Lexeme -> Senses).
// fetcher: функция сервиса, возвращающая плоский список всех найденных элементов.
// getKey: функция, извлекающая Foreign Key из элемента.
func buildGroupedBatcher[V any](
	fetcher func(context.Context, []uuid.UUID) ([]V, error),
	getKey func(V) uuid.UUID,
) func(context.Context, []uuid.UUID) ([][]V, []error) {
	return func(ctx context.Context, keys []uuid.UUID) ([][]V, []error) {
		items, err := fetcher(ctx, keys)
		if err != nil {
			errs := make([]error, len(keys))
			for i := range errs {
				errs[i] = err
			}
			return nil, errs
		}

		// Группируем по ключу
		grouped := make(map[uuid.UUID][]V)
		for _, item := range items {
			key := getKey(item)
			grouped[key] = append(grouped[key], item)
		}

		// Расставляем в порядке запроса
		result := make([][]V, len(keys))
		for i, key := range keys {
			result[i] = grouped[key]
		}
		return result, nil
	}
}

// buildOneToOneBatcher создает функцию для One-to-One связей (например, Card -> SRSState).
func buildOneToOneBatcher[V any](
	fetcher func(context.Context, []uuid.UUID) ([]V, error),
	getKey func(V) uuid.UUID,
) func(context.Context, []uuid.UUID) ([]*V, []error) {
	return func(ctx context.Context, keys []uuid.UUID) ([]*V, []error) {
		items, err := fetcher(ctx, keys)
		if err != nil {
			errs := make([]error, len(keys))
			for i := range errs {
				errs[i] = err
			}
			return nil, errs
		}

		// Мапим по ключу
		mapped := make(map[uuid.UUID]V)
		for _, item := range items {
			mapped[getKey(item)] = item
		}

		// Расставляем
		result := make([]*V, len(keys))
		for i, key := range keys {
			if item, ok := mapped[key]; ok {
				result[i] = &item // Берем адрес копии, так как V - структура
			} else {
				result[i] = nil
			}
		}
		return result, nil
	}
}

// --- Context Middleware ---

type ctxKey string

const loadersKey = ctxKey("dataloaders")

// Middleware инжектит лоадеры в контекст запроса.
func Middleware(svc *loader.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			loaders := NewLoaders(svc)
			ctx := context.WithValue(r.Context(), loadersKey, loaders)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// For возвращает лоадеры из контекста.
func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}
