package example

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/model"
)

// Create создаёт новый example.
func (r *Repo) Create(ctx context.Context, example *model.Example) error {
	if example == nil {
		return database.ErrInvalidInput
	}

	if example.MeaningID == 0 || example.SentenceEn == "" {
		return database.ErrInvalidInput
	}

	query, args, err := database.Builder.
		Insert(tableName).
		Columns("meaning_id", "sentence_en", "sentence_ru", "source_name").
		Values(
			example.MeaningID,
			example.SentenceEn,
			database.NullString(example.SentenceRu),
			example.SourceName,
		).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return err
	}

	err = r.q.QueryRowContext(ctx, query, args...).Scan(&example.ID)
	if err != nil {
		return database.WrapDBError(err)
	}

	return nil
}

// CreateBatch создаёт несколько examples за один запрос.
func (r *Repo) CreateBatch(ctx context.Context, examples []*model.Example) error {
	if len(examples) == 0 {
		return nil
	}

	qb := database.Builder.
		Insert(tableName).
		Columns("meaning_id", "sentence_en", "sentence_ru", "source_name")

	for _, ex := range examples {
		if ex == nil || ex.MeaningID == 0 || ex.SentenceEn == "" {
			return database.ErrInvalidInput
		}
		qb = qb.Values(
			ex.MeaningID,
			ex.SentenceEn,
			database.NullString(ex.SentenceRu),
			ex.SourceName,
		)
	}

	qb = qb.Suffix("RETURNING id")

	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	rows, err := r.q.QueryContext(ctx, query, args...)
	if err != nil {
		return database.WrapDBError(err)
	}
	defer rows.Close()

	i := 0
	for rows.Next() {
		if err := rows.Scan(&examples[i].ID); err != nil {
			return err
		}
		i++
	}

	return rows.Err()
}

// Update обновляет example.
func (r *Repo) Update(ctx context.Context, example *model.Example) error {
	if example == nil {
		return database.ErrInvalidInput
	}

	if example.SentenceEn == "" {
		return database.ErrInvalidInput
	}

	query, args, err := database.Builder.
		Update(tableName).
		Set("sentence_en", example.SentenceEn).
		Set("sentence_ru", database.NullString(example.SentenceRu)).
		Set("source_name", example.SourceName).
		Where(squirrel.Eq{"id": example.ID}).
		ToSql()
	if err != nil {
		return err
	}

	result, err := r.q.ExecContext(ctx, query, args...)
	if err != nil {
		return database.WrapDBError(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return database.ErrNotFound
	}

	return nil
}

// Delete удаляет example по ID.
func (r *Repo) Delete(ctx context.Context, id int64) error {
	query, args, err := database.Builder.
		Delete(tableName).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	result, err := r.q.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return database.ErrNotFound
	}

	return nil
}

// DeleteByMeaningID удаляет все examples для указанного meaning.
func (r *Repo) DeleteByMeaningID(ctx context.Context, meaningID int64) (int64, error) {
	query, args, err := database.Builder.
		Delete(tableName).
		Where(squirrel.Eq{"meaning_id": meaningID}).
		ToSql()
	if err != nil {
		return 0, err
	}

	result, err := r.q.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
