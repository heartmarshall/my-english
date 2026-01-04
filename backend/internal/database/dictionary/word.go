package dictionary

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

// GetByText возвращает список слов из словаря по тексту (из разных источников).
func (r *Repo) GetByText(ctx context.Context, text string) ([]model.DictionaryWord, error) {
	builder := database.Builder.
		Select(schema.DictionaryWords.All()...).
		From(schema.DictionaryWords.Name.String()).
		Where(schema.DictionaryWords.Text.Eq(text))

	return database.NewQuery[model.DictionaryWord](r.q, builder).List(ctx)
}

// SearchSimilar использует триграммный поиск.
func (r *Repo) SearchSimilar(ctx context.Context, query string, limit int, similarityThreshold float64) ([]model.DictionaryWord, error) {
	if limit <= 0 {
		limit = 10
	}
	if similarityThreshold <= 0 {
		similarityThreshold = 0.3
	}

	trigramCond := squirrel.Or{
		squirrel.Expr("word_similarity(?, ?) > ?", query, schema.DictionaryWords.Text, similarityThreshold),
		squirrel.Expr("? % ?", schema.DictionaryWords.Text, query),
	}

	// 1. Внутренний запрос с расчетом similarity
	innerBuilder := database.Builder.
		Select(schema.DictionaryWords.All()...).
		Column(squirrel.Expr("word_similarity(?, ?) AS similarity", query, schema.DictionaryWords.Text)).
		From(schema.DictionaryWords.Name.String()).
		Where(trigramCond)

	// 2. Внешний запрос для чистой проекции и сортировки
	finalCols := make([]string, 0)
	for _, col := range schema.DictionaryWords.All() {
		finalCols = append(finalCols, "sub."+schema.Column(col).Bare())
	}

	outerBuilder := database.Builder.
		Select(finalCols...).
		FromSelect(innerBuilder, "sub").
		OrderBy("sub.similarity DESC").
		Limit(uint64(limit))

	return database.NewQuery[model.DictionaryWord](r.q, outerBuilder).List(ctx)
}

// SaveWordData сохраняет слово и все его связи в БД (Write-through cache).
func (r *Repo) SaveWordData(ctx context.Context, data *model.DictionaryWordData) error {
	// Приводим интерфейс Querier к *pgxpool.Pool для управления транзакцией,
	// либо r.q уже должен уметь запускать транзакции (если используется TxManager).
	// В текущей реализации database.WithTx требует *pgxpool.Pool.
	pool, ok := r.q.(*pgxpool.Pool)
	if !ok {
		return fmt.Errorf("repository querier is not a connection pool, cannot start transaction")
	}

	return database.WithTx(ctx, pool, func(ctx context.Context, tx database.Querier) error {
		now := r.clock.Now()

		// 1. Сохраняем/Обновляем слово
		// ON CONFLICT (text, source) DO UPDATE ...
		builder := database.Builder.
			Insert(schema.DictionaryWords.Name.String()).
			Columns(schema.DictionaryWords.InsertColumns()...).
			Values(
				data.Word.Text,
				data.Word.Transcription,
				data.Word.AudioURL,
				data.Word.FrequencyRank,
				data.Word.Source,
				data.Word.SourceID,
				now, // created_at
				now, // updated_at
			).
			Suffix("ON CONFLICT (text, source) DO UPDATE SET updated_at = EXCLUDED.updated_at, transcription = EXCLUDED.transcription, audio_url = EXCLUDED.audio_url RETURNING id")

		id, err := database.ExecInsertWithReturn[int64](ctx, tx, builder)
		if err != nil {
			return fmt.Errorf("failed to upsert word: %w", err)
		}
		data.Word.ID = id

		// 2. Удаляем старые значения для этого слова (стратегия полной замены для кэша)
		delBuilder := database.Builder.
			Delete(schema.DictionaryMeanings.Name.String()).
			Where(schema.DictionaryMeanings.DictionaryWordID.Eq(id))

		if _, err := database.ExecOnly(ctx, tx, delBuilder); err != nil {
			return fmt.Errorf("failed to cleanup meanings: %w", err)
		}

		// 3. Вставляем новые значения
		for _, mData := range data.Meanings {
			mData.Meaning.DictionaryWordID = id
			mData.Meaning.CreatedAt = now
			mData.Meaning.UpdatedAt = now

			mBuilder := database.Builder.
				Insert(schema.DictionaryMeanings.Name.String()).
				Columns(schema.DictionaryMeanings.InsertColumns()...).
				Values(
					mData.Meaning.DictionaryWordID,
					mData.Meaning.PartOfSpeech,
					mData.Meaning.DefinitionEn,
					mData.Meaning.CefrLevel,
					mData.Meaning.ImageURL,
					mData.Meaning.OrderIndex,
					now, now,
				).
				Suffix("RETURNING id")

			mID, err := database.ExecInsertWithReturn[int64](ctx, tx, mBuilder)
			if err != nil {
				return fmt.Errorf("failed to insert meaning: %w", err)
			}

			// 4. Вставляем переводы (если есть)
			if len(mData.Translations) > 0 {
				tBuilder := database.Builder.
					Insert(schema.DictionaryTranslations.Name.String()).
					Columns(schema.DictionaryTranslations.InsertColumns()...)

				for _, t := range mData.Translations {
					tBuilder = tBuilder.Values(mID, t.TranslationRu, now)
				}

				tBuilder = tBuilder.Suffix("ON CONFLICT DO NOTHING")

				if _, err := database.ExecOnly(ctx, tx, tBuilder); err != nil {
					return fmt.Errorf("failed to insert translations: %w", err)
				}
			}
		}

		return nil
	})
}
