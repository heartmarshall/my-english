package example

import (
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/model"
	"github.com/jackc/pgx/v5"
)

func (r *Repo) scanRow(s database.Scanner) (*model.Example, error) {
	var (
		id         int64
		meaningID  int64
		sentenceEn string
		sentenceRu *string
		sourceName *string
	)

	err := s.Scan(&id, &meaningID, &sentenceEn, &sentenceRu, &sourceName)
	if err != nil {
		return nil, err
	}

	ex := &model.Example{
		ID:         id,
		MeaningID:  meaningID,
		SentenceEn: sentenceEn,
		SentenceRu: sentenceRu,
	}

	if sourceName != nil {
		src := model.ExampleSource(*sourceName)
		ex.SourceName = &src
	}

	return ex, nil
}

func (r *Repo) scanRows(rows pgx.Rows) ([]*model.Example, error) {
	examples := make([]*model.Example, 0)

	for rows.Next() {
		ex, err := r.scanRow(rows)
		if err != nil {
			return nil, err
		}
		examples = append(examples, ex)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return examples, nil
}
