package meaning

import (
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/Masterminds/squirrel"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/model"
)

// scanRow сканирует одну строку в model.Meaning.
func (r *Repo) scanRow(s database.Scanner) (*model.Meaning, error) {
	var (
		id             int64
		wordID         int64
		partOfSpeech   string
		definitionEn   *string
		translationRu  string
		cefrLevel      *string
		imageURL       *string
		learningStatus string
		nextReviewAt   *time.Time
		interval       *int
		easeFactor     *float64
		reviewCount    *int
		createdAt      *time.Time
		updatedAt      *time.Time
	)

	err := s.Scan(
		&id, &wordID, &partOfSpeech, &definitionEn, &translationRu,
		&cefrLevel, &imageURL, &learningStatus, &nextReviewAt,
		&interval, &easeFactor, &reviewCount, &createdAt, &updatedAt,
	)
	if err != nil {
		return nil, err
	}

	m := &model.Meaning{
		ID:             id,
		WordID:         wordID,
		PartOfSpeech:   model.PartOfSpeech(partOfSpeech),
		DefinitionEn:   definitionEn,
		TranslationRu:  translationRu,
		CefrLevel:      cefrLevel,
		ImageURL:       imageURL,
		LearningStatus: model.LearningStatus(learningStatus),
		NextReviewAt:   nextReviewAt,
		Interval:       interval,
		EaseFactor:     easeFactor,
		ReviewCount:    reviewCount,
	}

	if createdAt != nil {
		m.CreatedAt = *createdAt
	}
	if updatedAt != nil {
		m.UpdatedAt = *updatedAt
	}

	return m, nil
}

// scanRows сканирует несколько строк в слайс model.Meaning.
func (r *Repo) scanRows(rows pgx.Rows) ([]*model.Meaning, error) {
	meanings := make([]*model.Meaning, 0)

	for rows.Next() {
		meaning, err := r.scanRow(rows)
		if err != nil {
			return nil, err
		}
		meanings = append(meanings, meaning)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return meanings, nil
}

// applyFilter применяет фильтры к query builder.
func applyFilter(qb squirrel.SelectBuilder, filter *Filter) squirrel.SelectBuilder {
	if filter == nil {
		return qb
	}

	if filter.WordID != nil {
		qb = qb.Where(squirrel.Eq{"word_id": *filter.WordID})
	}
	if filter.PartOfSpeech != nil {
		qb = qb.Where(squirrel.Eq{"part_of_speech": *filter.PartOfSpeech})
	}
	if filter.LearningStatus != nil {
		qb = qb.Where(squirrel.Eq{"learning_status": *filter.LearningStatus})
	}

	return qb
}
