package word

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

func (r *Repo) GetByID(ctx context.Context, id int64) (model.Word, error) {
	query, args, err := database.Builder.
		Select(schema.Words.All()...).
		From(schema.Words.Name.String()).
		Where(schema.Words.ID.Eq(id)).
		ToSql()
	if err != nil {
		return model.Word{}, err
	}

	word, err := database.GetOne[model.Word](ctx, r.q, query, args...)
	if err != nil {
		return model.Word{}, err
	}
	return *word, nil
}

func (r *Repo) GetByText(ctx context.Context, text string) (model.Word, error) {
	query, args, err := database.Builder.
		Select(schema.Words.All()...).
		From(schema.Words.Name.String()).
		Where(schema.Words.Text.Eq(text)).
		ToSql()
	if err != nil {
		return model.Word{}, err
	}

	word, err := database.GetOne[model.Word](ctx, r.q, query, args...)
	if err != nil {
		return model.Word{}, err
	}
	return *word, nil
}

func (r *Repo) List(ctx context.Context, filter *model.WordFilter, limit, offset int) ([]model.Word, error) {
	limit, offset = database.NormalizePagination(limit, offset)

	qb := database.Builder.
		Select(schema.Words.All()...).
		From(schema.Words.Name.String())

	qb = applyFilter(qb, filter)

	qb = qb.
		OrderBy(schema.Words.CreatedAt.Desc()).
		Limit(uint64(limit)).
		Offset(uint64(offset))

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	return database.Select[model.Word](ctx, r.q, query, args...)
}

func (r *Repo) Count(ctx context.Context, filter *model.WordFilter) (int, error) {
	qb := database.Builder.Select("COUNT(*)").From(schema.Words.Name.String())
	qb = applyFilter(qb, filter)

	query, args, err := qb.ToSql()
	if err != nil {
		return 0, err
	}

	return database.GetScalar[int](ctx, r.q, query, args...)
}

func (r *Repo) Exists(ctx context.Context, id int64) (bool, error) {
	query, args, err := database.Builder.
		Select("1").
		From(schema.Words.Name.String()).
		Where(schema.Words.ID.Eq(id)).
		Limit(1).
		ToSql()
	if err != nil {
		return false, err
	}
	return database.CheckExists(ctx, r.q, query, args...)
}

func applyFilter(qb squirrel.SelectBuilder, filter *model.WordFilter) squirrel.SelectBuilder {
	if filter == nil {
		return qb
	}
	if filter.Search != nil && *filter.Search != "" {
		return qb.Where(schema.Words.Text.ILike("%" + *filter.Search + "%"))
	}
	return qb
}
