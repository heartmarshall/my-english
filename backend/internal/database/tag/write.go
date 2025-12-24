package tag

import (
	"context"
	"strings"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/schema"
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
		Insert(schema.Tags.Name.String()).
		Columns(schema.Tags.NameCol.Bare()).
		Values(name).
		Suffix("RETURNING " + schema.Tags.ID.Bare()).
		ToSql()
	if err != nil {
		return err
	}

	err = r.q.QueryRow(ctx, query, args...).Scan(&tag.ID)
	if err != nil {
		return database.WrapDBError(err)
	}

	tag.Name = name
	return nil
}

// GetOrCreate возвращает существующий tag или создаёт новый.
func (r *Repo) GetOrCreate(ctx context.Context, name string) (model.Tag, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return model.Tag{}, database.ErrInvalidInput
	}

	// Пробуем найти существующий
	tag, err := r.GetByName(ctx, name)
	if err == nil {
		return tag, nil
	}

	if err != database.ErrNotFound {
		return model.Tag{}, err
	}

	// Создаём новый
	tagPtr := &model.Tag{Name: name}
	if err := r.Create(ctx, tagPtr); err != nil {
		// Возможен race condition — проверяем ещё раз
		if database.IsDuplicateError(err) {
			return r.GetByName(ctx, name)
		}
		return model.Tag{}, err
	}

	return *tagPtr, nil
}

// Delete удаляет tag по ID.
func (r *Repo) Delete(ctx context.Context, id int64) error {
	query, args, err := database.Builder.
		Delete(schema.Tags.Name.String()).
		Where(schema.Tags.ID.Eq(id)).
		ToSql()
	if err != nil {
		return err
	}

	commandTag, err := r.q.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return database.ErrNotFound
	}

	return nil
}
