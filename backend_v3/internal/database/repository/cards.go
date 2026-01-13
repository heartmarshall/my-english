package repository

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

type CardRepository struct {
	*Base[model.Card]
}

func NewCardRepository(q database.Querier) *CardRepository {
	return &CardRepository{
		Base: NewBase[model.Card](q, schema.Cards.Name.String(), schema.Cards.Columns()),
	}
}

func (r *CardRepository) GetByEntryID(ctx context.Context, entryID uuid.UUID) (*model.Card, error) {
	return r.FindOneBy(ctx, schema.Cards.EntryID.String(), entryID)
}

// GetByIDForUpdate получает карточку и блокирует строку (SELECT FOR UPDATE)
func (r *CardRepository) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*model.Card, error) {
	query := r.SelectBuilder().
		Where(squirrel.Eq{schema.Cards.ID.String(): id}).
		Suffix("FOR UPDATE")

	return r.GetOne(ctx, query)
}

func (r *CardRepository) GetDueCards(ctx context.Context, now time.Time, limit int) ([]model.Card, error) {
	b := r.SelectBuilder().
		Where(squirrel.LtOrEq{schema.Cards.NextReviewAt.String(): now}).
		OrderBy(schema.Cards.NextReviewAt.String() + " ASC").
		Limit(uint64(limit))

	return r.List(ctx, b)
}

// DashboardStatsDTO содержит агрегированную статистику
type DashboardStatsDTO struct {
	TotalWords    int `db:"total_words"`
	TotalCards    int `db:"total_cards"`
	NewCards      int `db:"new_cards"`
	LearningCards int `db:"learning_cards"`
	ReviewCards   int `db:"review_cards"`
	MasteredCards int `db:"mastered_cards"`
	DueToday      int `db:"due_today"`
}

// GetDashboardStats возвращает статистику одним запросом
func (r *CardRepository) GetDashboardStats(ctx context.Context) (*DashboardStatsDTO, error) {
	sql := `
		SELECT
			(SELECT COUNT(*) FROM dictionary_entries)::int as total_words,
			COUNT(*)::int as total_cards,
			COUNT(*) FILTER (WHERE status = 'NEW')::int as new_cards,
			COUNT(*) FILTER (WHERE status = 'LEARNING')::int as learning_cards,
			COUNT(*) FILTER (WHERE status = 'REVIEW')::int as review_cards,
			COUNT(*) FILTER (WHERE status = 'MASTERED')::int as mastered_cards,
			COUNT(*) FILTER (WHERE next_review_at <= NOW())::int as due_today
		FROM cards
	`

	var stats DashboardStatsDTO
	if err := pgxscan.Get(ctx, r.Q(), &stats, sql); err != nil {
		return nil, database.WrapDBError(err)
	}
	return &stats, nil
}

// Create создает новую карточку с дефолтными значениями
func (r *CardRepository) Create(ctx context.Context, card *model.Card) (*model.Card, error) {
	// Устанавливаем дефолтные значения если не заданы
	if card.Status == "" {
		card.Status = model.StatusNew
	}
	if card.EaseFactor == 0 {
		card.EaseFactor = 2.5
	}

	insert := r.InsertBuilder().
		Columns(schema.Cards.InsertColumns()...).
		Values(
			card.EntryID,
			card.Status,
			card.NextReviewAt,
			card.IntervalDays,
			card.EaseFactor,
		)

	return r.InsertReturning(ctx, insert)
}

// Update обновляет карточку
func (r *CardRepository) Update(ctx context.Context, id uuid.UUID, card *model.Card) (*model.Card, error) {
	update := r.UpdateBuilder().
		Set("status", card.Status).
		Set("next_review_at", card.NextReviewAt).
		Set("interval_days", card.IntervalDays).
		Set("ease_factor", card.EaseFactor).
		Where(squirrel.Eq{schema.Cards.ID.String(): id})

	return r.Base.Update(ctx, update)
}

// UpdateSRSFields обновляет только SRS поля карточки
func (r *CardRepository) UpdateSRSFields(ctx context.Context, id uuid.UUID, status model.LearningStatus, nextReviewAt *time.Time, intervalDays int, easeFactor float64) error {
	update := r.UpdateBuilder().
		Set("status", status).
		Set("next_review_at", nextReviewAt).
		Set("interval_days", intervalDays).
		Set("ease_factor", easeFactor).
		Where(squirrel.Eq{schema.Cards.ID.String(): id})

	_, err := r.Base.Update(ctx, update)
	return err
}

// Delete удаляет карточку
func (r *CardRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.Base.Delete(ctx, schema.Cards.ID.String(), id)
}

type ReviewLogRepository struct {
	*Base[model.ReviewLog]
}

func NewReviewLogRepository(q database.Querier) *ReviewLogRepository {
	return &ReviewLogRepository{
		Base: NewBase[model.ReviewLog](q, schema.ReviewLogs.Name.String(), schema.ReviewLogs.Columns()),
	}
}

// Create создает запись о review
func (r *ReviewLogRepository) Create(ctx context.Context, log *model.ReviewLog) (*model.ReviewLog, error) {
	insert := r.InsertBuilder().
		Columns(schema.ReviewLogs.InsertColumns()...).
		Values(log.CardID, log.Grade, log.DurationMs)

	return r.InsertReturning(ctx, insert)
}

// ListByCardID возвращает историю повторений для карточки
func (r *ReviewLogRepository) ListByCardID(ctx context.Context, cardID uuid.UUID, limit int) ([]model.ReviewLog, error) {
	query := r.SelectBuilder().
		Where(squirrel.Eq{schema.ReviewLogs.CardID.String(): cardID}).
		OrderBy(schema.ReviewLogs.ReviewedAt.String() + " DESC")

	if limit > 0 {
		query = query.Limit(uint64(limit))
	}

	return r.List(ctx, query)
}
