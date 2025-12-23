package tag

import (
	"github.com/jackc/pgx/v5"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/model"
)

func (r *Repo) scanRow(s database.Scanner) (*model.Tag, error) {
	var (
		id   int64
		name string
	)

	err := s.Scan(&id, &name)
	if err != nil {
		return nil, err
	}

	return &model.Tag{
		ID:   id,
		Name: name,
	}, nil
}

func (r *Repo) scanRows(rows pgx.Rows) ([]*model.Tag, error) {
	tags := make([]*model.Tag, 0)

	for rows.Next() {
		tag, err := r.scanRow(rows)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}
