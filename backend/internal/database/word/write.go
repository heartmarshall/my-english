package word

import (
	"context"
	"strings"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

func (r *Repo) Create(ctx context.Context, word *model.Word) error {
	if word == nil || strings.TrimSpace(word.Text) == "" {
		return database.ErrInvalidInput
	}

	now := r.clock.Now()

	query, args, err := database.Builder.
		Insert(schema.Words.String()).
		Columns(
			schema.WordColumns.Text.String(),
			schema.WordColumns.Transcription.String(),
			schema.WordColumns.AudioURL.String(),
			schema.WordColumns.FrequencyRank.String(),
			schema.WordColumns.CreatedAt.String(),
		).
		Values(
			word.Text,
			word.Transcription,
			word.AudioURL,
			word.FrequencyRank,
			now,
		).
		Suffix(schema.WordColumns.ID.Returning()).
		ToSql()
	if err != nil {
		return err
	}

	// Для возврата ID используем стандартный Scan, так как это одно поле
	err = r.q.QueryRow(ctx, query, args...).Scan(&word.ID)
	if err != nil {
		return database.WrapDBError(err)
	}

	word.CreatedAt = now
	return nil
}

func (r *Repo) Update(ctx context.Context, word *model.Word) error {
	if word == nil || strings.TrimSpace(word.Text) == "" {
		return database.ErrInvalidInput
	}

	text := strings.TrimSpace(strings.ToLower(word.Text))

	query, args, err := database.Builder.
		Update(schema.Words.String()).
		Set(schema.WordColumns.Text.String(), text).
		Set(schema.WordColumns.Transcription.String(), word.Transcription). // Прямая передача
		Set(schema.WordColumns.AudioURL.String(), word.AudioURL).
		Set(schema.WordColumns.FrequencyRank.String(), word.FrequencyRank).
		Where(schema.WordColumns.ID.Eq(word.ID)).
		ToSql()
	if err != nil {
		return err
	}

	cmd, err := r.q.Exec(ctx, query, args...)
	if err != nil {
		return database.WrapDBError(err)
	}
	if cmd.RowsAffected() == 0 {
		return database.ErrNotFound
	}

	word.Text = text
	return nil
}

func (r *Repo) Delete(ctx context.Context, id int64) error {
	query, args, err := database.Builder.Delete(schema.Words.String()).Where(schema.WordColumns.ID.Eq(id)).ToSql()
	if err != nil {
		return err
	}
	cmd, err := r.q.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return database.ErrNotFound
	}
	return nil
}
