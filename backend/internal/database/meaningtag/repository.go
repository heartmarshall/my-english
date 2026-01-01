package meaningtag

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

type Repo struct {
	q database.Querier
}

func New(q database.Querier) *Repo {
	return &Repo{q: q}
}

func (r *Repo) AttachTag(ctx context.Context, meaningID, tagID int64) error {
	builder := database.Builder.
		Insert(schema.MeaningTags.Name.String()).
		Columns(
			schema.MeaningTags.MeaningID.Bare(),
			schema.MeaningTags.TagID.Bare(),
		).
		Values(meaningID, tagID).
		Suffix("ON CONFLICT DO NOTHING")

	_, err := database.ExecOnly(ctx, r.q, builder)
	return err
}

func (r *Repo) AttachTags(ctx context.Context, meaningID int64, tagIDs []int64) error {
	if len(tagIDs) == 0 {
		return nil
	}
	qb := database.Builder.
		Insert(schema.MeaningTags.Name.String()).
		Columns(
			schema.MeaningTags.MeaningID.Bare(),
			schema.MeaningTags.TagID.Bare(),
		)

	for _, tagID := range tagIDs {
		qb = qb.Values(meaningID, tagID)
	}
	qb = qb.Suffix("ON CONFLICT DO NOTHING")

	_, err := database.ExecOnly(ctx, r.q, qb)
	return err
}

func (r *Repo) DetachTag(ctx context.Context, meaningID, tagID int64) error {
	builder := database.Builder.
		Delete(schema.MeaningTags.Name.String()).
		Where(squirrel.And{
			schema.MeaningTags.MeaningID.Eq(meaningID),
			schema.MeaningTags.TagID.Eq(tagID),
		})

	_, err := database.ExecOnly(ctx, r.q, builder)
	return err
}

func (r *Repo) DetachAllFromMeaning(ctx context.Context, meaningID int64) error {
	builder := database.Builder.
		Delete(schema.MeaningTags.Name.String()).
		Where(schema.MeaningTags.MeaningID.Eq(meaningID))

	_, err := database.ExecOnly(ctx, r.q, builder)
	return err
}

// GetTagIDsByMeaningID возвращает ID тегов для meaning.
func (r *Repo) GetTagIDsByMeaningID(ctx context.Context, meaningID int64) ([]int64, error) {
	builder := database.Builder.
		Select(schema.MeaningTags.TagID.String()).
		From(schema.MeaningTags.Name.String()).
		Where(schema.MeaningTags.MeaningID.Eq(meaningID))

	return database.NewQuery[int64](r.q, builder).List(ctx)
}

// GetMeaningIDsByTagID возвращает ID meanings для tag.
func (r *Repo) GetMeaningIDsByTagID(ctx context.Context, tagID int64) ([]int64, error) {
	builder := database.Builder.
		Select(schema.MeaningTags.MeaningID.String()).
		From(schema.MeaningTags.Name.String()).
		Where(schema.MeaningTags.TagID.Eq(tagID))

	return database.NewQuery[int64](r.q, builder).List(ctx)
}

// GetByMeaningIDs возвращает все связи для нескольких meanings.
func (r *Repo) GetByMeaningIDs(ctx context.Context, meaningIDs []int64) ([]model.MeaningTag, error) {
	if len(meaningIDs) == 0 {
		return make([]model.MeaningTag, 0), nil
	}
	builder := database.Builder.
		Select(
			schema.MeaningTags.MeaningID.String(),
			schema.MeaningTags.TagID.String(),
		).
		From(schema.MeaningTags.Name.String()).
		Where(schema.MeaningTags.MeaningID.In(meaningIDs))

	return database.NewQuery[model.MeaningTag](r.q, builder).List(ctx)
}

// SyncTags синхронизирует теги для meaning, выполняя только необходимые операции.
// Вычисляет разницу между текущими и новыми тегами, чтобы избежать лишних операций.
func (r *Repo) SyncTags(ctx context.Context, meaningID int64, tagIDs []int64) error {
	// Получаем текущие теги
	currentTagIDs, err := r.GetTagIDsByMeaningID(ctx, meaningID)
	if err != nil {
		return err
	}

	// Вычисляем diff: какие теги нужно добавить, а какие удалить
	toAdd, toRemove := computeTagDiff(currentTagIDs, tagIDs)

	// Удаляем теги, которых больше нет
	if len(toRemove) > 0 {
		for _, tagID := range toRemove {
			if err := r.DetachTag(ctx, meaningID, tagID); err != nil {
				return err
			}
		}
	}

	// Добавляем новые теги
	if len(toAdd) > 0 {
		if err := r.AttachTags(ctx, meaningID, toAdd); err != nil {
			return err
		}
	}

	return nil
}

// computeTagDiff вычисляет разницу между текущими и новыми тегами.
// Возвращает слайсы тегов для добавления и удаления.
func computeTagDiff(current, desired []int64) (toAdd, toRemove []int64) {
	// Создаем мапы для быстрого поиска
	currentMap := make(map[int64]bool, len(current))
	for _, id := range current {
		currentMap[id] = true
	}

	desiredMap := make(map[int64]bool, len(desired))
	for _, id := range desired {
		desiredMap[id] = true
	}

	// Находим теги для добавления (есть в desired, но нет в current)
	toAdd = make([]int64, 0)
	for _, id := range desired {
		if !currentMap[id] {
			toAdd = append(toAdd, id)
		}
	}

	// Находим теги для удаления (есть в current, но нет в desired)
	toRemove = make([]int64, 0)
	for _, id := range current {
		if !desiredMap[id] {
			toRemove = append(toRemove, id)
		}
	}

	return toAdd, toRemove
}
