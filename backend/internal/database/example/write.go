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

	query, args, err := database.Builder.
		Insert(schema.Examples.String()).
		Columns(
			schema.ExampleColumns.MeaningID.String(),
			schema.ExampleColumns.SentenceEn.String(),
			schema.ExampleColumns.SentenceRu.String(),
			schema.ExampleColumns.SourceName.String(),
		).
		Values(
			example.MeaningID,
			example.SentenceEn,
			example.SentenceRu,
			example.SourceName,
		).
		Suffix(schema.ExampleColumns.ID.Returning()).
		ToSql()
	if err != nil {
		return err
	}

	err = r.q.QueryRow(ctx, query, args...).Scan(&example.ID)
	if err != nil {
		return database.WrapDBError(err)
	}

	return nil
}

// CreateBatch создаёт несколько examples за один запрос.
// Использует pgx.CollectRows для эффективного сбора ID.
func (r *Repo) CreateBatch(ctx context.Context, examples []*model.Example) error {
	if len(examples) == 0 {
		return nil
	}

	qb := database.Builder.
		Insert(schema.Examples.String()).
		Columns(
			schema.ExampleColumns.MeaningID.String(),
			schema.ExampleColumns.SentenceEn.String(),
			schema.ExampleColumns.SentenceRu.String(),
			schema.ExampleColumns.SourceName.String(),
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

	qb = qb.Suffix("RETURNING " + schema.ExampleColumns.ID.String())

	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	// Выполняем запрос
	rows, err := r.q.Query(ctx, query, args...)
	if err != nil {
		return database.WrapDBError(err)
	}

	// pgx.CollectRows автоматически закрывает rows и обрабатывает ошибки сканирования.
	// pgx.RowTo[int64] — эффективный маппер для одиночной колонки.
	ids, err := pgx.CollectRows(rows, pgx.RowTo[int64])
	if err != nil {
		return database.WrapDBError(err)
	}

	// Присваиваем полученные ID обратно в структуры
	// Порядок RETURNING в PostgreSQL соответствует порядку VALUES (для INSERT).
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

	query, args, err := database.Builder.
		Update(schema.Examples.String()).
		Set(schema.ExampleColumns.SentenceEn.String(), example.SentenceEn).
		Set(schema.ExampleColumns.SentenceRu.String(), example.SentenceRu).
		Set(schema.ExampleColumns.SourceName.String(), example.SourceName).
		Where(schema.ExampleColumns.ID.Eq(example.ID)).
		ToSql()
	if err != nil {
		return err
	}

	commandTag, err := r.q.Exec(ctx, query, args...)
	if err != nil {
		return database.WrapDBError(err)
	}

	if commandTag.RowsAffected() == 0 {
		return database.ErrNotFound
	}

	return nil
}

// Delete удаляет example по ID.
func (r *Repo) Delete(ctx context.Context, id int64) error {
	query, args, err := database.Builder.
		Delete(schema.Examples.String()).
		Where(schema.ExampleColumns.ID.Eq(id)).
		ToSql()
	if err != nil {
		return err
	}

	commandTag, err := r.q.Exec(ctx, query, args...)
	if err != nil {
		return database.WrapDBError(err)
	}

	if commandTag.RowsAffected() == 0 {
		return database.ErrNotFound
	}

	return nil
}

// DeleteByMeaningID удаляет все examples для указанного meaning.
func (r *Repo) DeleteByMeaningID(ctx context.Context, meaningID int64) (int64, error) {
	query, args, err := database.Builder.
		Delete(schema.Examples.String()).
		Where(schema.ExampleColumns.MeaningID.Eq(meaningID)).
		ToSql()
	if err != nil {
		return 0, err
	}

	commandTag, err := r.q.Exec(ctx, query, args...)
	if err != nil {
		return 0, database.WrapDBError(err)
	}

	return commandTag.RowsAffected(), nil
}
