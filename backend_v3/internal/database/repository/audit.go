package repository

import (
	"context"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

type AuditRepository struct {
	*Base[model.AuditRecord]
}

func NewAuditRepository(q database.Querier) *AuditRepository {
	return &AuditRepository{
		Base: NewBase[model.AuditRecord](q, schema.AuditRecords.Name.String(), schema.AuditRecords.Columns()),
	}
}

// Create создает новую запись аудита
func (r *AuditRepository) Create(ctx context.Context, audit *model.AuditRecord) (*model.AuditRecord, error) {
	insert := r.InsertBuilder().
		Columns(schema.AuditRecords.InsertColumns()...).
		Values(
			audit.EntityType,
			audit.EntityID,
			audit.Action,
			audit.Changes,
		)

	return r.InsertReturning(ctx, insert)
}
