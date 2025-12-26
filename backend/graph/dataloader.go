package graph

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/heartmarshall/my-english/internal/model"
	"github.com/vikstrous/dataloadgen"
)

// ctxKey — ключ для хранения loaders в context.
type ctxKey string

const loadersKey ctxKey = "dataloaders"

// Loaders содержит все DataLoaders.
type Loaders struct {
	ExamplesByMeaningID *dataloadgen.Loader[int64, *[]*model.Example]
	TagsByMeaningID     *dataloadgen.Loader[int64, *[]*model.Tag]
	MeaningsByWordID    *dataloadgen.Loader[int64, *[]*model.Meaning]
}

// LoaderService определяет интерфейс для batch-загрузки данных.
// Реализуется сервисом loader.Service.
type LoaderService interface {
	GetMeaningsByWordIDs(ctx context.Context, wordIDs []int64) ([]model.Meaning, error)
	GetExamplesByMeaningIDs(ctx context.Context, meaningIDs []int64) ([]model.Example, error)
	GetTagsByMeaningIDs(ctx context.Context, meaningIDs []int64) ([]model.MeaningTag, error)
	GetTagsByIDs(ctx context.Context, ids []int64) ([]model.Tag, error)
}

// LoaderDeps — зависимости для DataLoaders.
type LoaderDeps struct {
	Loader LoaderService
}

// NewLoaders создаёт новые DataLoaders.
func NewLoaders(deps LoaderDeps) *Loaders {
	return &Loaders{
		ExamplesByMeaningID: dataloadgen.NewLoader(
			newExampleBatchFunc(deps.Loader),
			dataloadgen.WithWait(2*time.Millisecond),
		),
		TagsByMeaningID: dataloadgen.NewLoader(
			newTagBatchFunc(deps.Loader),
			dataloadgen.WithWait(2*time.Millisecond),
		),
		MeaningsByWordID: dataloadgen.NewLoader(
			newMeaningBatchFunc(deps.Loader),
			dataloadgen.WithWait(2*time.Millisecond),
		),
	}
}

// newMeaningBatchFunc создаёт batch функцию для загрузки meanings.
func newMeaningBatchFunc(loader LoaderService) func(ctx context.Context, wordIDs []int64) ([]*[]*model.Meaning, []error) {
	return func(ctx context.Context, wordIDs []int64) ([]*[]*model.Meaning, []error) {
		meanings, err := loader.GetMeaningsByWordIDs(ctx, wordIDs)
		if err != nil {
			errs := make([]error, len(wordIDs))
			for i := range errs {
				errs[i] = err
			}
			return nil, errs
		}

		// Группируем по wordID
		grouped := make(map[int64][]*model.Meaning)
		for i := range meanings {
			m := &meanings[i]
			grouped[m.WordID] = append(grouped[m.WordID], m)
		}

		// Формируем результат
		result := make([]*[]*model.Meaning, len(wordIDs))
		for i, id := range wordIDs {
			ms := grouped[id]
			if ms == nil {
				ms = make([]*model.Meaning, 0)
			}
			result[i] = &ms
		}

		return result, nil
	}
}

// newExampleBatchFunc создаёт batch функцию для загрузки examples.
func newExampleBatchFunc(loader LoaderService) func(ctx context.Context, meaningIDs []int64) ([]*[]*model.Example, []error) {
	return func(ctx context.Context, meaningIDs []int64) ([]*[]*model.Example, []error) {
		// Загружаем все examples одним запросом
		examples, err := loader.GetExamplesByMeaningIDs(ctx, meaningIDs)
		if err != nil {
			// Логируем ошибку, но возвращаем пустые массивы вместо ошибки
			// чтобы не ломать весь запрос, если examples отсутствуют
			slog.Warn("Failed to load examples", "error", err, "meaningIDs", meaningIDs)
			// Возвращаем пустые массивы для всех ключей
			result := make([]*[]*model.Example, len(meaningIDs))
			for i := range result {
				empty := make([]*model.Example, 0)
				result[i] = &empty
			}
			return result, nil
		}

		// Группируем по meaningID
		grouped := make(map[int64][]*model.Example)
		for i := range examples {
			ex := &examples[i]
			grouped[ex.MeaningID] = append(grouped[ex.MeaningID], ex)
		}

		// Формируем результат в том же порядке, что и входные ключи
		result := make([]*[]*model.Example, len(meaningIDs))
		for i, id := range meaningIDs {
			exs := grouped[id]
			if exs == nil {
				exs = make([]*model.Example, 0)
			}
			result[i] = &exs
		}

		return result, nil
	}
}

// newTagBatchFunc создаёт batch функцию для загрузки tags.
func newTagBatchFunc(loader LoaderService) func(ctx context.Context, meaningIDs []int64) ([]*[]*model.Tag, []error) {
	return func(ctx context.Context, meaningIDs []int64) ([]*[]*model.Tag, []error) {
		// Загружаем связи meaning-tag
		meaningTags, err := loader.GetTagsByMeaningIDs(ctx, meaningIDs)
		if err != nil {
			errs := make([]error, len(meaningIDs))
			for i := range errs {
				errs[i] = err
			}
			return nil, errs
		}

		// Собираем уникальные tagIDs
		tagIDSet := make(map[int64]struct{})
		for i := range meaningTags {
			tagIDSet[meaningTags[i].TagID] = struct{}{}
		}

		tagIDs := make([]int64, 0, len(tagIDSet))
		for id := range tagIDSet {
			tagIDs = append(tagIDs, id)
		}

		// Загружаем теги
		var tagMap map[int64]*model.Tag
		if len(tagIDs) > 0 {
			tags, err := loader.GetTagsByIDs(ctx, tagIDs)
			if err != nil {
				errs := make([]error, len(meaningIDs))
				for i := range errs {
					errs[i] = err
				}
				return nil, errs
			}

			tagMap = make(map[int64]*model.Tag, len(tags))
			for i := range tags {
				t := &tags[i]
				tagMap[t.ID] = t
			}
		} else {
			tagMap = make(map[int64]*model.Tag)
		}

		// Группируем теги по meaningID
		grouped := make(map[int64][]*model.Tag)
		for i := range meaningTags {
			mt := &meaningTags[i]
			if tag, ok := tagMap[mt.TagID]; ok {
				grouped[mt.MeaningID] = append(grouped[mt.MeaningID], tag)
			}
		}

		// Формируем результат
		result := make([]*[]*model.Tag, len(meaningIDs))
		for i, id := range meaningIDs {
			tags := grouped[id]
			if tags == nil {
				tags = make([]*model.Tag, 0)
			}
			result[i] = &tags
		}

		return result, nil
	}
}

// DataLoaderMiddleware добавляет DataLoaders в context для каждого запроса.
func DataLoaderMiddleware(deps LoaderDeps) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Создаём новые loaders для каждого запроса (важно для батчинга)
			loaders := NewLoaders(deps)
			ctx := context.WithValue(r.Context(), loadersKey, loaders)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetLoaders извлекает Loaders из context.
func GetLoaders(ctx context.Context) *Loaders {
	loaders, ok := ctx.Value(loadersKey).(*Loaders)
	if !ok {
		return nil
	}
	return loaders
}

// LoadExamplesForMeaning загружает examples для meaning через DataLoader.
func LoadExamplesForMeaning(ctx context.Context, meaningID int64) ([]*model.Example, error) {
	loaders := GetLoaders(ctx)
	if loaders == nil {
		return make([]*model.Example, 0), nil
	}
	examples, err := loaders.ExamplesByMeaningID.Load(ctx, meaningID)
	if err != nil {
		return nil, err
	}
	if examples == nil {
		return make([]*model.Example, 0), nil
	}
	return *examples, nil
}

// LoadTagsForMeaning загружает tags для meaning через DataLoader.
func LoadTagsForMeaning(ctx context.Context, meaningID int64) ([]*model.Tag, error) {
	loaders := GetLoaders(ctx)
	if loaders == nil {
		return make([]*model.Tag, 0), nil
	}
	tags, err := loaders.TagsByMeaningID.Load(ctx, meaningID)
	if err != nil {
		return nil, err
	}
	if tags == nil {
		return make([]*model.Tag, 0), nil
	}
	return *tags, nil
}

// LoadMeaningsForWord загружает meanings для word через DataLoader.
func LoadMeaningsForWord(ctx context.Context, wordID int64) ([]*model.Meaning, error) {
	loaders := GetLoaders(ctx)
	if loaders == nil {
		return make([]*model.Meaning, 0), nil
	}
	meanings, err := loaders.MeaningsByWordID.Load(ctx, wordID)
	if err != nil {
		return nil, err
	}
	if meanings == nil {
		return make([]*model.Meaning, 0), nil
	}
	return *meanings, nil
}
