package study

import (
	"context"
	"time"

	"github.com/heartmarshall/my-english/internal/model"
	"github.com/heartmarshall/my-english/internal/service/study/srs"
)

// MeaningRepository определяет интерфейс для работы со значениями.
type MeaningRepository interface {
	GetByID(ctx context.Context, id int64) (model.Meaning, error)
	GetStudyQueue(ctx context.Context, limit int) ([]model.Meaning, error)
	GetStats(ctx context.Context) (*model.Stats, error)
}

// MeaningSRSRepository определяет интерфейс для обновления SRS данных.
type MeaningSRSRepository interface {
	UpdateSRS(ctx context.Context, id int64, update *srs.Update) error
}

// Clock — интерфейс для получения времени (для тестирования).
type Clock interface {
	Now() time.Time
}

// RealClock — реальная реализация Clock.
type RealClock struct{}

// Now возвращает текущее время.
func (RealClock) Now() time.Time {
	return time.Now()
}
