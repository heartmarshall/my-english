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

	builder := database.Builder.
		Insert(schema.Words.Name.String()).
		Columns(schema.Words.InsertColumns()...).
		Values(
			word.Text,
			word.Transcription,
			word.AudioURL,
			word.FrequencyRank,
			now,
		).
		Suffix("RETURNING " + schema.Words.ID.Bare())

	id, err := database.ExecInsertWithReturn[int64](ctx, r.q, builder)
	if err != nil {
		return err
	}

	word.ID = id

	word.CreatedAt = now
	return nil
}

func (r *Repo) Update(ctx context.Context, word *model.Word) error {
	if word == nil || strings.TrimSpace(word.Text) == "" {
		return database.ErrInvalidInput
	}

	text := strings.TrimSpace(strings.ToLower(word.Text))

	builder := database.Builder.
		Update(schema.Words.Name.String()).
		Set(schema.Words.Text.Bare(), text).
		Set(schema.Words.Transcription.Bare(), word.Transcription). // Прямая передача
		Set(schema.Words.AudioURL.Bare(), word.AudioURL).
		Set(schema.Words.FrequencyRank.Bare(), word.FrequencyRank).
		Where(schema.Words.ID.Eq(word.ID))

	rowsAffected, err := database.ExecOnly(ctx, r.q, builder)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return database.ErrNotFound
	}

	word.Text = text
	return nil
}

func (r *Repo) Delete(ctx context.Context, id int64) error {
	builder := database.Builder.Delete(schema.Words.Name.String()).Where(schema.Words.ID.Eq(id))

	rowsAffected, err := database.ExecOnly(ctx, r.q, builder)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return database.ErrNotFound
	}
	return nil
}
