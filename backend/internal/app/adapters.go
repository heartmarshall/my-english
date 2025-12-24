package app

import (
	"context"

	"github.com/heartmarshall/my-english/internal/database/meaning"
	"github.com/heartmarshall/my-english/internal/service/study/srs"
)

// --- SRS Adapter ---

// SRSAdapter адаптирует meaning.Repo для study.MeaningSRSRepository.
type SRSAdapter struct {
	repo *meaning.Repo
}

// NewSRSAdapter создаёт адаптер.
func NewSRSAdapter(repo *meaning.Repo) *SRSAdapter {
	return &SRSAdapter{repo: repo}
}

// UpdateSRS реализует study.MeaningSRSRepository.
func (a *SRSAdapter) UpdateSRS(ctx context.Context, id int64, update *srs.Update) error {
	// Конвертируем srs.Update в meaning.SRSUpdate
	meaningUpdate := &meaning.SRSUpdate{
		LearningStatus: update.LearningStatus,
		NextReviewAt:   update.NextReviewAt,
		Interval:       update.Interval,
		EaseFactor:     update.EaseFactor,
		ReviewCount:    update.ReviewCount,
	}
	return a.repo.UpdateSRS(ctx, id, meaningUpdate)
}
