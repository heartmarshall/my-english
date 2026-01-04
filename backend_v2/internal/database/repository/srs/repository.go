package srs

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/repository"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

type Repository struct {
	*repository.Base[model.SRSState]
}

func New(q database.Querier) *Repository {
	return &Repository{
		Base: repository.NewBase[model.SRSState](q, schema.SRSStates.Name.String(), schema.SRSStates.Columns()),
	}
}

// GetByCardID возвращает состояние SRS для карточки.
func (r *Repository) GetByCardID(ctx context.Context, cardID uuid.UUID) (*model.SRSState, error) {
	// В этой таблице card_id является Primary Key
	return r.Base.GetByID(ctx, schema.SRSStates.CardID.String(), cardID)
}

// GetByCardIDForUpdate возвращает состояние SRS с блокировкой строки (FOR UPDATE).
// Используется для предотвращения race conditions при конкурентных review запросах.
func (r *Repository) GetByCardIDForUpdate(ctx context.Context, cardID uuid.UUID) (*model.SRSState, error) {
	query := r.SelectBuilder().
		Where(schema.SRSStates.CardID.Eq(cardID)).
		Suffix("FOR UPDATE")
	return r.GetOne(ctx, query)
}

// ListDueForReview возвращает карточки, которые нужно повторить (due_date < now).
func (r *Repository) ListDueForReview(ctx context.Context, limit int) ([]model.SRSState, error) {
	return r.ListDueForReviewWithFilter(ctx, nil, limit)
}

// ListDueForReviewWithFilter возвращает карточки для повторения с фильтрацией по статусам.
func (r *Repository) ListDueForReviewWithFilter(ctx context.Context, statuses []model.LearningStatus, limit int) ([]model.SRSState, error) {
	query := r.SelectBuilder().
		Where(schema.SRSStates.DueDate.LtOrEq(time.Now())).
		OrderBy(schema.SRSStates.DueDate.Asc())

	// Если указаны статусы, фильтруем по ним, иначе исключаем только MASTERED
	if len(statuses) > 0 {
		query = query.Where(schema.SRSStates.Status.In(statuses))
	} else {
		query = query.Where(schema.SRSStates.Status.NotEq(model.LearningStatusMastered))
	}

	query = repository.ApplyOptions(query, repository.WithLimit(limit))
	return r.List(ctx, query)
}

// ListDueForReviewWithTags возвращает карточки для повторения с фильтрацией по статусам и тегам.
// Фильтрация по тегам выполняется на уровне SQL для оптимальной производительности.
// Карточка должна иметь ВСЕ указанные теги (AND логика).
func (r *Repository) ListDueForReviewWithTags(
	ctx context.Context,
	statuses []model.LearningStatus,
	tagNames []string,
	limit int,
) ([]model.SRSState, error) {
	query := r.SelectBuilder().
		Where(schema.SRSStates.DueDate.LtOrEq(time.Now())).
		OrderBy(schema.SRSStates.DueDate.Asc())

	// Фильтр по статусам
	if len(statuses) > 0 {
		query = query.Where(schema.SRSStates.Status.In(statuses))
	} else {
		query = query.Where(schema.SRSStates.Status.NotEq(model.LearningStatusMastered))
	}

	// Фильтр по тегам через подзапрос
	// Карточка должна иметь все указанные теги
	if len(tagNames) > 0 {
		// Подзапрос: выбираем card_id, которые имеют все нужные теги
		tagSubquery := repository.Builder.
			Select(schema.CardTags.CardID.String()).
			From(schema.CardTags.Name.String()).
			Join(schema.Tags.Name.String()+" ON "+schema.Tags.ID.Qualified()+" = "+schema.CardTags.TagID.Qualified()).
			Where(schema.Tags.NameCol.In(tagNames)).
			GroupBy(schema.CardTags.CardID.String()).
			Having("COUNT(DISTINCT "+schema.Tags.ID.Qualified()+") = ?", len(tagNames))

		subSQL, subArgs, err := tagSubquery.ToSql()
		if err != nil {
			return nil, err
		}

		// Добавляем условие: card_id должен быть в результате подзапроса
		query = query.Where(schema.SRSStates.CardID.String()+" IN ("+subSQL+")", subArgs...)
	}

	query = repository.ApplyOptions(query, repository.WithLimit(limit))
	return r.List(ctx, query)
}

// Upsert создает или обновляет состояние SRS.
func (r *Repository) Upsert(ctx context.Context, state *model.SRSState) (*model.SRSState, error) {
	insert := r.InsertBuilder().
		Columns(schema.SRSStates.InsertColumns()...).
		Values(
			state.CardID,
			state.Status,
			state.DueDate,
			state.AlgorithmData, // pgx/v5 автоматически сериализует map[string]any в JSONB
			state.LastReviewAt,
		).
		Suffix("ON CONFLICT (card_id) DO UPDATE SET status = EXCLUDED.status, due_date = EXCLUDED.due_date, algorithm_data = EXCLUDED.algorithm_data, last_review_at = EXCLUDED.last_review_at RETURNING *")

	return r.InsertReturning(ctx, insert)
}

func (r *Repository) ListByCardIDs(ctx context.Context, cardIDs []uuid.UUID) ([]model.SRSState, error) {
	return r.FindBy(ctx, schema.SRSStates.CardID.String(), cardIDs)
}
