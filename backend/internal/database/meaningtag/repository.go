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
	query, args, err := database.Builder.
		Insert(schema.MeaningTags.Name.String()).
		Columns(
			schema.MeaningTags.MeaningID.Bare(),
			schema.MeaningTags.TagID.Bare(),
		).
		Values(meaningID, tagID).
		Suffix("ON CONFLICT DO NOTHING").
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.q.Exec(ctx, query, args...)
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

	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}
	_, err = r.q.Exec(ctx, query, args...)
	return err
}

func (r *Repo) DetachTag(ctx context.Context, meaningID, tagID int64) error {
	query, args, err := database.Builder.
		Delete(schema.MeaningTags.Name.String()).
		Where(squirrel.And{
			schema.MeaningTags.MeaningID.Eq(meaningID),
			schema.MeaningTags.TagID.Eq(tagID),
		}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.q.Exec(ctx, query, args...)
	return err
}

func (r *Repo) DetachAllFromMeaning(ctx context.Context, meaningID int64) error {
	query, args, err := database.Builder.
		Delete(schema.MeaningTags.Name.String()).
		Where(schema.MeaningTags.MeaningID.Eq(meaningID)).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.q.Exec(ctx, query, args...)
	return err
}

// GetTagIDsByMeaningID возвращает ID тегов для meaning.
// Используем SelectScalars для получения []int64.
func (r *Repo) GetTagIDsByMeaningID(ctx context.Context, meaningID int64) ([]int64, error) {
	query, args, err := database.Builder.
		Select(schema.MeaningTags.TagID.String()).
		From(schema.MeaningTags.Name.String()).
		Where(schema.MeaningTags.MeaningID.Eq(meaningID)).
		ToSql()
	if err != nil {
		return nil, err
	}
	return database.SelectScalars[int64](ctx, r.q, query, args...)
}

// GetMeaningIDsByTagID возвращает ID meanings для tag.
func (r *Repo) GetMeaningIDsByTagID(ctx context.Context, tagID int64) ([]int64, error) {
	query, args, err := database.Builder.
		Select(schema.MeaningTags.MeaningID.String()).
		From(schema.MeaningTags.Name.String()).
		Where(schema.MeaningTags.TagID.Eq(tagID)).
		ToSql()
	if err != nil {
		return nil, err
	}
	return database.SelectScalars[int64](ctx, r.q, query, args...)
}

// GetByMeaningIDs возвращает все связи для нескольких meanings.
// Используем Select для получения списка структур.
func (r *Repo) GetByMeaningIDs(ctx context.Context, meaningIDs []int64) ([]model.MeaningTag, error) {
	if len(meaningIDs) == 0 {
		return make([]model.MeaningTag, 0), nil
	}
	query, args, err := database.Builder.
		Select(
			schema.MeaningTags.MeaningID.String(),
			schema.MeaningTags.TagID.String(),
		).
		From(schema.MeaningTags.Name.String()).
		Where(schema.MeaningTags.MeaningID.In(meaningIDs)).
		ToSql()
	if err != nil {
		return nil, err
	}
	return database.Select[model.MeaningTag](ctx, r.q, query, args...)
}

func (r *Repo) SyncTags(ctx context.Context, meaningID int64, tagIDs []int64) error {
	// TODO: Оптимизировать через вычисление Diff (insert/delete), чтобы избежать bloat таблицы
	if err := r.DetachAllFromMeaning(ctx, meaningID); err != nil {
		return err
	}
	return r.AttachTags(ctx, meaningID, tagIDs)
}
