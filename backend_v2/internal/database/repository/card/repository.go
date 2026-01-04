package card

import (
	"context"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/repository"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

// escapeLikePattern экранирует спецсимволы LIKE паттерна для безопасного поиска.
func escapeLikePattern(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "%", "\\%")
	s = strings.ReplaceAll(s, "_", "\\_")
	return s
}

// Repository работает с таблицей cards.
type Repository struct {
	*repository.Base[model.Card]
}

// New создаёт новый репозиторий.
func New(q database.Querier) *Repository {
	return &Repository{
		Base: repository.NewBase[model.Card](q, schema.Cards.Name.String(), schema.Cards.Columns()),
	}
}

// GetByID возвращает активную (не удалённую) карточку по ID.
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*model.Card, error) {
	query := r.SelectBuilder().
		Where(schema.Cards.ID.Eq(id)).
		Where(schema.Cards.IsDeleted.Eq(false))

	return r.GetOne(ctx, query)
}

// ListByIDs возвращает список активных (не удалённых) карточек по списку ID.
func (r *Repository) ListByIDs(ctx context.Context, ids []uuid.UUID) ([]model.Card, error) {
	if len(ids) == 0 {
		return []model.Card{}, nil
	}

	query := r.SelectBuilder().
		Where(schema.Cards.ID.In(ids)).
		Where(schema.Cards.IsDeleted.Eq(false))

	return r.List(ctx, query)
}

// ListActive возвращает список активных карточек с возможностью фильтрации и пагинации.
func (r *Repository) ListActive(ctx context.Context, opts ...repository.QueryOption) ([]model.Card, error) {
	// Базовое условие - только не удаленные
	query := r.SelectBuilder().
		Where(schema.Cards.IsDeleted.Eq(false))

	// Применяем внешние опции (сортировка, лимит, офсет)
	query = repository.ApplyOptions(query, opts...)

	return r.List(ctx, query)
}

// ListWithFilters возвращает список карточек с фильтрацией по тегам, статусам и поиску.
func (r *Repository) ListWithFilters(
	ctx context.Context,
	tagNames []string,
	statuses []model.LearningStatus,
	searchText *string,
	opts ...repository.QueryOption,
) ([]model.Card, error) {
	query := r.SelectBuilder().
		Where(schema.Cards.IsDeleted.Eq(false))

	// Фильтрация по тегам (используем subquery)
	if len(tagNames) > 0 {
		// Subquery: находим card_id, которые имеют все указанные теги
		tagSubquery := repository.Builder.
			Select("DISTINCT "+schema.CardTags.CardID.String()).
			From(schema.CardTags.Name.String()).
			Join(schema.Tags.Name.String()+" ON "+schema.Tags.ID.Qualified()+" = "+schema.CardTags.TagID.Qualified()).
			Where(schema.Tags.NameCol.In(tagNames)).
			GroupBy(schema.CardTags.CardID.String()).
			Having("COUNT(DISTINCT "+schema.Tags.ID.Qualified()+") = ?", len(tagNames))

		sql, args, err := tagSubquery.ToSql()
		if err != nil {
			return nil, err
		}

		// Используем IN с subquery через squirrel
		query = query.Where(squirrel.Expr(schema.Cards.ID.Qualified()+" IN ("+sql+")", args...))
	}

	// Фильтрация по статусам (JOIN с srs_states)
	if len(statuses) > 0 {
		query = query.
			Join(schema.SRSStates.Name.String() + " ON " + schema.Cards.ID.Qualified() + " = " + schema.SRSStates.CardID.Qualified()).
			Where(schema.SRSStates.Status.In(statuses))
	}

	// Поиск по тексту (custom_text)
	if searchText != nil && *searchText != "" {
		// Экранируем спецсимволы LIKE для безопасности
		escaped := escapeLikePattern(*searchText)
		searchPattern := "%" + escaped + "%"
		// Поиск в custom_text (ILIKE для case-insensitive поиска)
		query = query.Where(schema.Cards.CustomText.ILike(searchPattern))
	}

	// Если был JOIN, нужно добавить DISTINCT, чтобы избежать дубликатов
	// ВАЖНО: DISTINCT должен быть ДО LIMIT/OFFSET
	if len(statuses) > 0 {
		query = query.Distinct()
	}

	// Применяем внешние опции (сортировка, лимит, офсет) в конце
	query = repository.ApplyOptions(query, opts...)

	return r.List(ctx, query)
}

// Create создаёт новую карточку.
func (r *Repository) Create(ctx context.Context, card *model.Card) (*model.Card, error) {
	insert := r.InsertBuilder().
		Columns(schema.Cards.InsertColumns()...).
		Values(
			card.SenseID,
			card.CustomText,
			card.CustomTranscription,
			card.CustomTranslations,
			card.CustomNote,
			card.CustomImageURL,
		)

	return r.InsertReturning(ctx, insert)
}

// Update обновляет поля карточки.
func (r *Repository) Update(ctx context.Context, card *model.Card) (*model.Card, error) {
	now := time.Now()

	update := r.UpdateBuilder().
		Set(schema.Cards.CustomText.Bare(), card.CustomText).
		Set(schema.Cards.CustomTranscription.Bare(), card.CustomTranscription).
		Set(schema.Cards.CustomTranslations.Bare(), card.CustomTranslations).
		Set(schema.Cards.CustomNote.Bare(), card.CustomNote).
		Set(schema.Cards.CustomImageURL.Bare(), card.CustomImageURL).
		Set(schema.Cards.UpdatedAt.Bare(), now).
		Where(schema.Cards.ID.Eq(card.ID)).
		Where(schema.Cards.IsDeleted.Eq(false)) // Защита от обновления удаленных

	return r.UpdateReturning(ctx, update)
}

// SoftDelete помечает карточку как удалённую.
func (r *Repository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	update := r.UpdateBuilder().
		Set(schema.Cards.IsDeleted.Bare(), true).
		Set(schema.Cards.UpdatedAt.Bare(), time.Now()).
		Where(schema.Cards.ID.Eq(id)).
		Where(schema.Cards.IsDeleted.Eq(false)) // Если уже удалена, rowsAffected будет 0

	// Используем базовый Update (который возвращает int64), а не UpdateReturning
	affected, err := r.Base.Update(ctx, update)
	if err != nil {
		return err
	}
	if affected == 0 {
		return database.ErrNotFound
	}
	return nil
}
