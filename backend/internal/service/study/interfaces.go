package study

import (
	"context"
	"time"

	"github.com/heartmarshall/my-english/internal/model"
)

// MeaningRepository определяет интерфейс для работы со значениями.
type MeaningRepository interface {
	GetByID(ctx context.Context, id int64) (*model.Meaning, error)
	GetStudyQueue(ctx context.Context, limit int) ([]*model.Meaning, error)
	GetStats(ctx context.Context) (*model.Stats, error)
}

// MeaningSRSRepository определяет интерфейс для обновления SRS данных.
type MeaningSRSRepository interface {
	UpdateSRS(ctx context.Context, id int64, srs *SRSUpdate) error
}

// SRSUpdate содержит поля для обновления SRS данных.
type SRSUpdate struct {
	LearningStatus model.LearningStatus
	NextReviewAt   *time.Time
	Interval       *int
	EaseFactor     *float64
	ReviewCount    *int
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
