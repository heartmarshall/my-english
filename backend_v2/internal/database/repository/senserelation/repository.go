package senserelation

import (
	"context"

	"github.com/google/uuid"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/repository"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

type Repository struct {
	*repository.Base[model.SenseRelation]
}

func New(q database.Querier) *Repository {
	return &Repository{
		Base: repository.NewBase[model.SenseRelation](q, schema.SenseRelations.Name.String(), schema.SenseRelations.Columns()),
	}
}

// ListBySourceSenseID возвращает связи, где данный sense является источником.
func (r *Repository) ListBySourceSenseID(ctx context.Context, senseID uuid.UUID) ([]model.SenseRelation, error) {
	return r.FindBy(ctx, schema.SenseRelations.SourceSenseID.String(), senseID)
}

// ListByTargetSenseID возвращает связи, где данный sense является целью.
// Используется для двусторонних связей.
func (r *Repository) ListByTargetSenseID(ctx context.Context, senseID uuid.UUID) ([]model.SenseRelation, error) {
	return r.FindBy(ctx, schema.SenseRelations.TargetSenseID.String(), senseID)
}

// ListBySenseID возвращает все связи для смысла (как источник, так и цель).
func (r *Repository) ListBySenseID(ctx context.Context, senseID uuid.UUID) ([]model.SenseRelation, error) {
	// Получаем связи, где sense является источником
	sourceRelations, err := r.ListBySourceSenseID(ctx, senseID)
	if err != nil {
		return nil, err
	}

	// Получаем связи, где sense является целью (для двусторонних связей)
	targetRelations, err := r.ListByTargetSenseID(ctx, senseID)
	if err != nil {
		return nil, err
	}

	// Объединяем результаты
	result := make([]model.SenseRelation, 0, len(sourceRelations)+len(targetRelations))
	result = append(result, sourceRelations...)

	// Для целевых связей создаем обратные связи, если они двусторонние
	for _, rel := range targetRelations {
		if rel.IsBidirectional {
			// Создаем обратную связь для отображения
			reverseRel := model.SenseRelation{
				SourceSenseID:   rel.TargetSenseID,
				TargetSenseID:   rel.SourceSenseID,
				Type:            rel.Type,
				IsBidirectional: rel.IsBidirectional,
				SourceID:        rel.SourceID,
			}
			result = append(result, reverseRel)
		}
	}

	return result, nil
}

// ListBySenseIDs batch-загрузка связей для списка смыслов.
func (r *Repository) ListBySenseIDs(ctx context.Context, senseIDs []uuid.UUID) ([]model.SenseRelation, error) {
	if len(senseIDs) == 0 {
		return []model.SenseRelation{}, nil
	}

	// Получаем связи, где sense является источником
	sourceRelations, err := r.FindBy(ctx, schema.SenseRelations.SourceSenseID.String(), senseIDs)
	if err != nil {
		return nil, err
	}

	// Получаем связи, где sense является целью
	targetRelations, err := r.FindBy(ctx, schema.SenseRelations.TargetSenseID.String(), senseIDs)
	if err != nil {
		return nil, err
	}

	// Объединяем и обрабатываем двусторонние связи
	result := make([]model.SenseRelation, 0, len(sourceRelations)+len(targetRelations))
	result = append(result, sourceRelations...)

	// Для целевых связей создаем обратные связи, если они двусторонние
	for _, rel := range targetRelations {
		if rel.IsBidirectional {
			reverseRel := model.SenseRelation{
				SourceSenseID:   rel.TargetSenseID,
				TargetSenseID:   rel.SourceSenseID,
				Type:            rel.Type,
				IsBidirectional: rel.IsBidirectional,
				SourceID:        rel.SourceID,
			}
			result = append(result, reverseRel)
		}
	}

	return result, nil
}

// Create создаёт новую связь между смыслами.
func (r *Repository) Create(ctx context.Context, rel *model.SenseRelation) (*model.SenseRelation, error) {
	insert := r.InsertBuilder().
		Columns(schema.SenseRelations.InsertColumns()...).
		Values(
			rel.SourceSenseID,
			rel.TargetSenseID,
			rel.Type,
			rel.IsBidirectional,
			rel.SourceID,
		)

	return r.InsertReturning(ctx, insert)
}
