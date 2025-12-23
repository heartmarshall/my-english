package study

import (
	"context"

	"github.com/heartmarshall/my-english/internal/model"
)

// GetStats возвращает статистику изучения.
func (s *Service) GetStats(ctx context.Context) (*model.Stats, error) {
	return s.meanings.GetStats(ctx)
}
