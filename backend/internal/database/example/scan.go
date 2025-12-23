package example

import (
	"database/sql"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/model"
)

func (r *Repo) scanRow(s database.Scanner) (*model.Example, error) {
	var (
		id         int64
		meaningID  int64
		sentenceEn string
		sentenceRu sql.NullString
		sourceName sql.NullString
	)

	err := s.Scan(&id, &meaningID, &sentenceEn, &sentenceRu, &sourceName)
	if err != nil {
		return nil, err
	}

	ex := &model.Example{
		ID:         id,
		MeaningID:  meaningID,
		SentenceEn: sentenceEn,
		SentenceRu: database.PtrString(sentenceRu),
	}

	if sourceName.Valid {
		src := model.ExampleSource(sourceName.String)
		ex.SourceName = &src
	}

	return ex, nil
}

func (r *Repo) scanRows(rows *sql.Rows) ([]*model.Example, error) {
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
