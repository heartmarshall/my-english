package app

import (
	"context"

	"github.com/heartmarshall/my-english/internal/database/meaning"
	"github.com/heartmarshall/my-english/internal/service/study"
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
func (a *SRSAdapter) UpdateSRS(ctx context.Context, id int64, srs *study.SRSUpdate) error {
	// Конвертируем study.SRSUpdate в meaning.SRSUpdate
	meaningUpdate := &meaning.SRSUpdate{
		LearningStatus: srs.LearningStatus,
		NextReviewAt:   srs.NextReviewAt,
		Interval:       srs.Interval,
		EaseFactor:     srs.EaseFactor,
		ReviewCount:    srs.ReviewCount,
	}
	return a.repo.UpdateSRS(ctx, id, meaningUpdate)
}
