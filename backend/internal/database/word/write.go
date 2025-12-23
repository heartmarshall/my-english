package word

import (
	"context"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/model"
)

// Create создаёт новое слово в базе данных.
// Возвращает database.ErrInvalidInput если word == nil или word.Text пустой.
// Возвращает database.ErrDuplicate если слово уже существует.
func (r *Repo) Create(ctx context.Context, word *model.Word) error {
	if word == nil {
		return database.ErrInvalidInput
	}

	text := strings.TrimSpace(strings.ToLower(word.Text))
	if text == "" {
		return database.ErrInvalidInput
	}

	now := r.clock.Now()

	query, args, err := database.Builder.
		Insert(tableName).
		Columns("text", "transcription", "audio_url", "frequency_rank", "created_at").
		Values(
			text,
			database.NullString(word.Transcription),
			database.NullString(word.AudioURL),
			database.NullInt(word.FrequencyRank),
			now,
		).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return err
	}

	err = r.q.QueryRowContext(ctx, query, args...).Scan(&word.ID)
	if err != nil {
		return database.WrapDBError(err)
	}

	word.Text = text
	word.CreatedAt = now
	return nil
}

// Update обновляет существующее слово.
// Возвращает database.ErrInvalidInput если word == nil.
// Возвращает database.ErrNotFound, если слово не найдено.
// Возвращает database.ErrDuplicate если новый текст уже существует.
func (r *Repo) Update(ctx context.Context, word *model.Word) error {
	if word == nil {
		return database.ErrInvalidInput
	}

	text := strings.TrimSpace(strings.ToLower(word.Text))
	if text == "" {
		return database.ErrInvalidInput
	}

	query, args, err := database.Builder.
		Update(tableName).
		Set("text", text).
		Set("transcription", database.NullString(word.Transcription)).
		Set("audio_url", database.NullString(word.AudioURL)).
		Set("frequency_rank", database.NullInt(word.FrequencyRank)).
		Where(squirrel.Eq{"id": word.ID}).
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

	word.Text = text
	return nil
}

// Delete удаляет слово по ID.
// Возвращает database.ErrNotFound, если слово не найдено.
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
