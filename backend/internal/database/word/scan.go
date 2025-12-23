package word

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/model"
)

// scanRow сканирует одну строку в model.Word.
func (r *Repo) scanRow(s database.Scanner) (*model.Word, error) {
	var (
		id            int64
		text          string
		transcription sql.NullString
		audioURL      sql.NullString
		frequencyRank sql.NullInt64
		createdAt     sql.NullTime
	)

	err := s.Scan(&id, &text, &transcription, &audioURL, &frequencyRank, &createdAt)
	if err != nil {
		return nil, err
	}

	word := &model.Word{
		ID:            id,
		Text:          text,
		Transcription: database.PtrString(transcription),
		AudioURL:      database.PtrString(audioURL),
		FrequencyRank: database.PtrInt(frequencyRank),
	}

	if createdAt.Valid {
		word.CreatedAt = createdAt.Time
	}

	return word, nil
}

// scanRows сканирует несколько строк в слайс model.Word.
func (r *Repo) scanRows(rows *sql.Rows) ([]*model.Word, error) {
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
