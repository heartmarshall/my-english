package tag

import (
	"context"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/model"
)

// Create создаёт новый tag.
func (r *Repo) Create(ctx context.Context, tag *model.Tag) error {
	if tag == nil {
		return database.ErrInvalidInput
	}

	name := strings.TrimSpace(tag.Name)
	if name == "" {
		return database.ErrInvalidInput
	}

	query, args, err := database.Builder.
		Insert(tableName).
		Columns("name").
		Values(name).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return err
	}

	err = r.q.QueryRowContext(ctx, query, args...).Scan(&tag.ID)
	if err != nil {
		return database.WrapDBError(err)
	}

	tag.Name = name
	return nil
}

// GetOrCreate возвращает существующий tag или создаёт новый.
func (r *Repo) GetOrCreate(ctx context.Context, name string) (*model.Tag, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, database.ErrInvalidInput
	}

	// Пробуем найти существующий
	tag, err := r.GetByName(ctx, name)
	if err == nil {
		return tag, nil
	}

	if err != database.ErrNotFound {
		return nil, err
	}

	// Создаём новый
	tag = &model.Tag{Name: name}
	if err := r.Create(ctx, tag); err != nil {
		// Возможен race condition — проверяем ещё раз
		if database.IsDuplicateError(err) {
			return r.GetByName(ctx, name)
		}
		return nil, err
	}

	return tag, nil
}

// Delete удаляет tag по ID.
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
