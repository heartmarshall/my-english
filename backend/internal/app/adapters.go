package app

import (
	"context"

	"github.com/heartmarshall/my-english/internal/database/meaning"
	"github.com/heartmarshall/my-english/internal/database/meaningtag"
	"github.com/heartmarshall/my-english/internal/database/tag"
	"github.com/heartmarshall/my-english/internal/model"
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

// --- Tag Loader ---

// TagLoaderAdapter комбинирует tag и meaningtag репозитории.
type TagLoaderAdapter struct {
	tags       *tag.Repo
	meaningTag *meaningtag.Repo
}

// NewTagLoaderAdapter создаёт адаптер.
func NewTagLoaderAdapter(tags *tag.Repo, meaningTag *meaningtag.Repo) *TagLoaderAdapter {
	return &TagLoaderAdapter{
		tags:       tags,
		meaningTag: meaningTag,
	}
}

// GetByMeaningIDs возвращает связи meaning-tag.
func (a *TagLoaderAdapter) GetByMeaningIDs(ctx context.Context, meaningIDs []int64) ([]*model.MeaningTag, error) {
	return a.meaningTag.GetByMeaningIDs(ctx, meaningIDs)
}

// GetByIDs возвращает теги по их ID.
func (a *TagLoaderAdapter) GetByIDs(ctx context.Context, ids []int64) ([]*model.Tag, error) {
	return a.tags.GetByIDs(ctx, ids)
}
