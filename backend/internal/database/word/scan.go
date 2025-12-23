package word

import (
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/model"
	"github.com/jackc/pgx/v5"
)

// scanRow сканирует одну строку в model.Word.
func (r *Repo) scanRow(s database.Scanner) (*model.Word, error) {
	var (
		id            int64
		text          string
		transcription *string
		audioURL      *string
		frequencyRank *int64
		createdAt     *time.Time
	)

	err := s.Scan(&id, &text, &transcription, &audioURL, &frequencyRank, &createdAt)
	if err != nil {
		return nil, err
	}

	var freqRank *int
	if frequencyRank != nil {
		val := int(*frequencyRank)
		freqRank = &val
	}

	word := &model.Word{
		ID:            id,
		Text:          text,
		Transcription: transcription,
		AudioURL:      audioURL,
		FrequencyRank: freqRank,
	}

	if createdAt != nil {
		word.CreatedAt = *createdAt
	}

	return word, nil
}

// scanRows сканирует несколько строк в слайс model.Word.
func (r *Repo) scanRows(rows pgx.Rows) ([]*model.Word, error) {
	words := make([]*model.Word, 0)

	for rows.Next() {
		word, err := r.scanRow(rows)
		if err != nil {
			return nil, err
		}
		words = append(words, word)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return words, nil
}

// applyFilter применяет фильтры к query builder.
func applyFilter(qb squirrel.SelectBuilder, filter *model.WordFilter) squirrel.SelectBuilder {
	if filter == nil {
		return qb
	}

	if filter.Search != nil && *filter.Search != "" {
		qb = qb.Where(squirrel.ILike{"text": "%" + *filter.Search + "%"})
	}

	return qb
}
