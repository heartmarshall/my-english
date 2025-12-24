package study_test

import (
	"context"
	"time"

	"github.com/heartmarshall/my-english/internal/model"
	"github.com/heartmarshall/my-english/internal/service/study/srs"
)

// --- Mock Clock ---

type mockClock struct {
	now time.Time
}

func (c *mockClock) Now() time.Time {
	return c.now
}

// --- Mock Repositories ---

type mockMeaningRepository struct {
	GetByIDFunc       func(ctx context.Context, id int64) (model.Meaning, error)
	GetStudyQueueFunc func(ctx context.Context, limit int) ([]model.Meaning, error)
	GetStatsFunc      func(ctx context.Context) (*model.Stats, error)
}

func (m *mockMeaningRepository) GetByID(ctx context.Context, id int64) (model.Meaning, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return model.Meaning{}, nil
}

func (m *mockMeaningRepository) GetStudyQueue(ctx context.Context, limit int) ([]model.Meaning, error) {
	if m.GetStudyQueueFunc != nil {
		return m.GetStudyQueueFunc(ctx, limit)
	}
	return []model.Meaning{}, nil
}

func (m *mockMeaningRepository) GetStats(ctx context.Context) (*model.Stats, error) {
	if m.GetStatsFunc != nil {
		return m.GetStatsFunc(ctx)
	}
	return &model.Stats{}, nil
}

type mockSRSRepository struct {
	UpdateSRSFunc func(ctx context.Context, id int64, update *srs.Update) error
}

func (m *mockSRSRepository) UpdateSRS(ctx context.Context, id int64, update *srs.Update) error {
	if m.UpdateSRSFunc != nil {
		return m.UpdateSRSFunc(ctx, id, update)
	}
	return nil
}
