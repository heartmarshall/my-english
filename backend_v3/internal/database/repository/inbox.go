package repository

import (
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

type InboxRepository struct {
	*Base[model.InboxItem]
}

func NewInboxRepository(q database.Querier) *InboxRepository {
	return &InboxRepository{
		Base: NewBase[model.InboxItem](q, schema.InboxItems.Name.String(), schema.InboxItems.Columns()),
	}
}
