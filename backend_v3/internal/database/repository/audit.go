package repository

import (
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
