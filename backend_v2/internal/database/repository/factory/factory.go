package repository

import (
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/repository/card"
	"github.com/heartmarshall/my-english/internal/database/repository/cardtag"
	"github.com/heartmarshall/my-english/internal/database/repository/datasource"
	"github.com/heartmarshall/my-english/internal/database/repository/example"
	"github.com/heartmarshall/my-english/internal/database/repository/inbox"
	"github.com/heartmarshall/my-english/internal/database/repository/inflection"
	"github.com/heartmarshall/my-english/internal/database/repository/lexeme"
	"github.com/heartmarshall/my-english/internal/database/repository/pronunciation"
	"github.com/heartmarshall/my-english/internal/database/repository/review"
	"github.com/heartmarshall/my-english/internal/database/repository/sense"
	"github.com/heartmarshall/my-english/internal/database/repository/senserelation"
	"github.com/heartmarshall/my-english/internal/database/repository/sensetranslation"
	"github.com/heartmarshall/my-english/internal/database/repository/srs"
	"github.com/heartmarshall/my-english/internal/database/repository/tag"
)

// Factory позволяет создавать инстансы репозиториев с заданным экзекьютором (пул или транзакция).
type Factory struct{}

func NewFactory() *Factory {
	return &Factory{}
}

// --- System ---
func (f *Factory) DataSource(q database.Querier) *datasource.Repository {
	return datasource.New(q)
}

// --- Linguistic ---
func (f *Factory) Lexeme(q database.Querier) *lexeme.Repository {
	return lexeme.New(q)
}

func (f *Factory) Sense(q database.Querier) *sense.Repository {
	return sense.New(q)
}

func (f *Factory) SenseTranslation(q database.Querier) *sensetranslation.Repository {
	return sensetranslation.New(q)
}

func (f *Factory) SenseRelation(q database.Querier) *senserelation.Repository {
	return senserelation.New(q)
}

func (f *Factory) Inflection(q database.Querier) *inflection.Repository {
	return inflection.New(q)
}

func (f *Factory) Example(q database.Querier) *example.Repository {
	return example.New(q)
}

func (f *Factory) Pronunciation(q database.Querier) *pronunciation.Repository {
	return pronunciation.New(q)
}

// --- User & Study ---
func (f *Factory) Card(q database.Querier) *card.Repository {
	return card.New(q)
}

func (f *Factory) SRS(q database.Querier) *srs.Repository {
	return srs.New(q)
}

func (f *Factory) Review(q database.Querier) *review.Repository {
	return review.New(q)
}

func (f *Factory) Inbox(q database.Querier) *inbox.Repository {
	return inbox.New(q)
}

func (f *Factory) Tag(q database.Querier) *tag.Repository {
	return tag.New(q)
}

func (f *Factory) CardTag(q database.Querier) *cardtag.Repository {
	return cardtag.New(q)
}
