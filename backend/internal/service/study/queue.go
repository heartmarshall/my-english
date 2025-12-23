package study

import (
	"context"

	"github.com/heartmarshall/my-english/internal/model"
)

// GetStudyQueue возвращает очередь слов для изучения/повторения.
// Включает: новые слова (status=NEW) и слова для повторения (next_review_at < NOW).
func (s *Service) GetStudyQueue(ctx context.Context, limit int) ([]*model.Meaning, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return s.meanings.GetStudyQueue(ctx, limit)
}
