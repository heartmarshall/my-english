package tag

import (
	"context"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/repository"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

type Repository struct {
	*repository.Base[model.Tag]
}

func New(q database.Querier) *Repository {
	return &Repository{
		Base: repository.NewBase[model.Tag](q, schema.Tags.Name.String(), schema.Tags.Columns()),
	}
}

func (r *Repository) GetByName(ctx context.Context, name string) (*model.Tag, error) {
	return r.FindOneBy(ctx, schema.Tags.NameCol.String(), name)
}

func (r *Repository) Create(ctx context.Context, tag *model.Tag) (*model.Tag, error) {
	insert := r.InsertBuilder().
		Columns(schema.Tags.InsertColumns()...).
		Values(tag.Name, tag.ColorHex)

	return r.InsertReturning(ctx, insert)
}

// GetOrCreate атомарно возвращает существующий тег или создаёт новый.
// Использует ON CONFLICT DO UPDATE для обхода проблемы aborted transaction в PostgreSQL.
func (r *Repository) GetOrCreate(ctx context.Context, name string) (*model.Tag, error) {
	insert := r.InsertBuilder().
		Columns(schema.Tags.InsertColumns()...).
		Values(name, nil). // ColorHex = nil для новых тегов
		Suffix("ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING *")

	return r.InsertReturning(ctx, insert)
}

// GetOrCreateBatch атомарно создаёт или возвращает теги по списку имён.
// Возвращает слайс тегов в порядке, соответствующем входным именам (без дубликатов).
func (r *Repository) GetOrCreateBatch(ctx context.Context, names []string) ([]model.Tag, error) {
	if len(names) == 0 {
		return []model.Tag{}, nil
	}

	// Дедуплицируем входные имена
	seen := make(map[string]bool)
	uniqueNames := make([]string, 0, len(names))
	for _, name := range names {
		if name != "" && !seen[name] {
			seen[name] = true
			uniqueNames = append(uniqueNames, name)
		}
	}

	if len(uniqueNames) == 0 {
		return []model.Tag{}, nil
	}

	// Используем batch insert с ON CONFLICT
	insert := r.InsertBuilder().
		Columns(schema.Tags.InsertColumns()...)

	for _, name := range uniqueNames {
		insert = insert.Values(name, nil)
	}

	insert = insert.Suffix("ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING *")

	return r.InsertReturningMany(ctx, insert)
}

func (r *Repository) GetByIDs(ctx context.Context, ids []int) ([]model.Tag, error) {
	if len(ids) == 0 {
		return []model.Tag{}, nil
	}
	return r.FindBy(ctx, schema.Tags.ID.String(), ids)
}

// GetByNames возвращает теги по списку имён (batch запрос).
func (r *Repository) GetByNames(ctx context.Context, names []string) ([]model.Tag, error) {
	if len(names) == 0 {
		return []model.Tag{}, nil
	}
	return r.FindBy(ctx, schema.Tags.NameCol.String(), names)
}
