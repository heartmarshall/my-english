package dictionary

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

// GetRelationsByMeaningID возвращает все связи (синонимы и антонимы) для значения по его ID.
// Возвращает связи, где значение является как meaning_id_1, так и meaning_id_2.
func (r *Repo) GetRelationsByMeaningID(ctx context.Context, meaningID int64) ([]model.DictionarySynonymAntonym, error) {
	// Ищем связи, где значение является meaning_id_1 или meaning_id_2
	builder := database.Builder.
		Select(schema.DictionarySynonymsAntonyms.All()...).
		From(schema.DictionarySynonymsAntonyms.Name.String()).
		Where(squirrel.Or{
			schema.DictionarySynonymsAntonyms.MeaningID1.Eq(meaningID),
			schema.DictionarySynonymsAntonyms.MeaningID2.Eq(meaningID),
		}).
		OrderBy(schema.DictionarySynonymsAntonyms.RelationType.Asc(), schema.DictionarySynonymsAntonyms.CreatedAt.Asc())

	return database.NewQuery[model.DictionarySynonymAntonym](r.q, builder).List(ctx)
}

// GetRelationsByMeaningIDs возвращает связи для нескольких значений (batch loading).
func (r *Repo) GetRelationsByMeaningIDs(ctx context.Context, meaningIDs []int64) ([]model.DictionarySynonymAntonym, error) {
	if len(meaningIDs) == 0 {
		return []model.DictionarySynonymAntonym{}, nil
	}

	builder := database.Builder.
		Select(schema.DictionarySynonymsAntonyms.All()...).
		From(schema.DictionarySynonymsAntonyms.Name.String()).
		Where(squirrel.Or{
			schema.DictionarySynonymsAntonyms.MeaningID1.In(meaningIDs),
			schema.DictionarySynonymsAntonyms.MeaningID2.In(meaningIDs),
		}).
		OrderBy(schema.DictionarySynonymsAntonyms.MeaningID1.Asc(), schema.DictionarySynonymsAntonyms.RelationType.Asc())

	return database.NewQuery[model.DictionarySynonymAntonym](r.q, builder).List(ctx)
}

// GetSynonymsByMeaningID возвращает только синонимы для значения.
func (r *Repo) GetSynonymsByMeaningID(ctx context.Context, meaningID int64) ([]model.DictionarySynonymAntonym, error) {
	builder := database.Builder.
		Select(schema.DictionarySynonymsAntonyms.All()...).
		From(schema.DictionarySynonymsAntonyms.Name.String()).
		Where(squirrel.And{
			schema.DictionarySynonymsAntonyms.RelationType.Eq(string(model.RelationTypeSynonym)),
			squirrel.Or{
				schema.DictionarySynonymsAntonyms.MeaningID1.Eq(meaningID),
				schema.DictionarySynonymsAntonyms.MeaningID2.Eq(meaningID),
			},
		}).
		OrderBy(schema.DictionarySynonymsAntonyms.CreatedAt.Asc())

	return database.NewQuery[model.DictionarySynonymAntonym](r.q, builder).List(ctx)
}

// GetAntonymsByMeaningID возвращает только антонимы для значения.
func (r *Repo) GetAntonymsByMeaningID(ctx context.Context, meaningID int64) ([]model.DictionarySynonymAntonym, error) {
	builder := database.Builder.
		Select(schema.DictionarySynonymsAntonyms.All()...).
		From(schema.DictionarySynonymsAntonyms.Name.String()).
		Where(squirrel.And{
			schema.DictionarySynonymsAntonyms.RelationType.Eq(string(model.RelationTypeAntonym)),
			squirrel.Or{
				schema.DictionarySynonymsAntonyms.MeaningID1.Eq(meaningID),
				schema.DictionarySynonymsAntonyms.MeaningID2.Eq(meaningID),
			},
		}).
		OrderBy(schema.DictionarySynonymsAntonyms.CreatedAt.Asc())

	return database.NewQuery[model.DictionarySynonymAntonym](r.q, builder).List(ctx)
}

// CreateRelation создаёт новую связь между значениями.
// Автоматически упорядочивает meaning_id_1 и meaning_id_2 (меньший ID всегда в meaning_id_1).
func (r *Repo) CreateRelation(ctx context.Context, relation *model.DictionarySynonymAntonym) error {
	if relation == nil {
		return database.ErrInvalidInput
	}

	// Убеждаемся, что meaning_id_1 < meaning_id_2
	meaningID1 := relation.MeaningID1
	meaningID2 := relation.MeaningID2
	if meaningID1 > meaningID2 {
		meaningID1, meaningID2 = meaningID2, meaningID1
	}

	now := r.clock.Now()
	relation.CreatedAt = now
	relation.UpdatedAt = now

	builder := database.Builder.
		Insert(schema.DictionarySynonymsAntonyms.Name.String()).
		Columns(schema.DictionarySynonymsAntonyms.InsertColumns()...).
		Values(
			meaningID1,
			meaningID2,
			relation.RelationType,
			relation.CreatedAt,
			relation.UpdatedAt,
		).
		Suffix("ON CONFLICT (meaning_id_1, meaning_id_2, relation_type) DO UPDATE SET updated_at = EXCLUDED.updated_at RETURNING " + schema.DictionarySynonymsAntonyms.ID.Bare())

	id, err := database.ExecInsertWithReturn[int64](ctx, r.q, builder)
	if err != nil {
		return err
	}

	relation.ID = id
	relation.MeaningID1 = meaningID1
	relation.MeaningID2 = meaningID2
	return nil
}

// DeleteRelation удаляет связь по ID.
func (r *Repo) DeleteRelation(ctx context.Context, relationID int64) error {
	builder := database.Builder.
		Delete(schema.DictionarySynonymsAntonyms.Name.String()).
		Where(schema.DictionarySynonymsAntonyms.ID.Eq(relationID))

	rowsAffected, err := database.ExecOnly(ctx, r.q, builder)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return database.ErrNotFound
	}
	return nil
}

