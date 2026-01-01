package example

import (
	"context"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
	"github.com/jackc/pgx/v5"
)

// Create создаёт новый example.
func (r *Repo) Create(ctx context.Context, example *model.Example) error {
	if example == nil || example.MeaningID == 0 || example.SentenceEn == "" {
		return database.ErrInvalidInput
	}

	builder := database.Builder.
		Insert(schema.Examples.Name.String()).
		Columns(
			schema.Examples.MeaningID.Bare(),
			schema.Examples.SentenceEn.Bare(),
			schema.Examples.SentenceRu.Bare(),
			schema.Examples.SourceName.Bare(),
		).
		Values(
			example.MeaningID,
			example.SentenceEn,
			example.SentenceRu,
			example.SourceName,
		).
		Suffix("RETURNING " + schema.Examples.ID.Bare())

	id, err := database.ExecInsertWithReturn[int64](ctx, r.q, builder)
	if err != nil {
		return err
	}

	example.ID = id
	return nil
}

// CreateBatch создаёт несколько examples за один запрос.
// Использует pgx.CollectRows для эффективного сбора ID.
func (r *Repo) CreateBatch(ctx context.Context, examples []*model.Example) error {
	if len(examples) == 0 {
		return nil
	}

	qb := database.Builder.
		Insert(schema.Examples.Name.String()).
		Columns(
			schema.Examples.MeaningID.Bare(),
			schema.Examples.SentenceEn.Bare(),
			schema.Examples.SentenceRu.Bare(),
			schema.Examples.SourceName.Bare(),
		)

	for _, ex := range examples {
		if ex == nil || ex.MeaningID == 0 || ex.SentenceEn == "" {
			return database.ErrInvalidInput
		}
		qb = qb.Values(
			ex.MeaningID,
			ex.SentenceEn,
			ex.SentenceRu,
			ex.SourceName,
		)
	}

	qb = qb.Suffix("RETURNING " + schema.Examples.ID.Bare())

	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	// Выполняем запрос
	rows, err := r.q.Query(ctx, query, args...)
	if err != nil {
		return database.WrapDBError(err)
	}

	ids, err := pgx.CollectRows(rows, pgx.RowTo[int64])
	if err != nil {
		return database.WrapDBError(err)
	}

	for i, id := range ids {
		if i < len(examples) {
			examples[i].ID = id
		}
	}

	return nil
}

// Update обновляет example.
func (r *Repo) Update(ctx context.Context, example *model.Example) error {
	if example == nil || example.SentenceEn == "" {
		return database.ErrInvalidInput
	}

	builder := database.Builder.
		Update(schema.Examples.Name.String()).
		Set(schema.Examples.SentenceEn.Bare(), example.SentenceEn).
		Set(schema.Examples.SentenceRu.Bare(), example.SentenceRu).
		Set(schema.Examples.SourceName.Bare(), example.SourceName).
		Where(schema.Examples.ID.Eq(example.ID))

	rowsAffected, err := database.ExecOnly(ctx, r.q, builder)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return database.ErrNotFound
	}

	return nil
}

// Delete удаляет example по ID.
func (r *Repo) Delete(ctx context.Context, id int64) error {
	builder := database.Builder.
		Delete(schema.Examples.Name.String()).
		Where(schema.Examples.ID.Eq(id))

	rowsAffected, err := database.ExecOnly(ctx, r.q, builder)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return database.ErrNotFound
	}

	return nil
}

// DeleteByMeaningID удаляет все examples для указанного meaning.
func (r *Repo) DeleteByMeaningID(ctx context.Context, meaningID int64) (int64, error) {
	builder := database.Builder.
		Delete(schema.Examples.Name.String()).
		Where(schema.Examples.MeaningID.Eq(meaningID))

	return database.ExecOnly(ctx, r.q, builder)
}
