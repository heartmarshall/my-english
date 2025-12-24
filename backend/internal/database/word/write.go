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
		Insert(schema.Words.Name.String()).
		Columns(schema.Words.InsertColumns()...).
		Values(
			word.Text,
			word.Transcription,
			word.AudioURL,
			word.FrequencyRank,
			now,
		).
		Suffix("RETURNING " + schema.Words.ID.Bare()).
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
		Update(schema.Words.Name.String()).
		Set(schema.Words.Text.Bare(), text).
		Set(schema.Words.Transcription.Bare(), word.Transcription). // Прямая передача
		Set(schema.Words.AudioURL.Bare(), word.AudioURL).
		Set(schema.Words.FrequencyRank.Bare(), word.FrequencyRank).
		Where(schema.Words.ID.Eq(word.ID)).
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
	query, args, err := database.Builder.Delete(schema.Words.Name.String()).Where(schema.Words.ID.Eq(id)).ToSql()
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
