package app

import (
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/example"
	"github.com/heartmarshall/my-english/internal/database/meaning"
	"github.com/heartmarshall/my-english/internal/database/meaningtag"
	"github.com/heartmarshall/my-english/internal/database/tag"
	"github.com/heartmarshall/my-english/internal/database/translation"
	"github.com/heartmarshall/my-english/internal/database/word"
	wordservice "github.com/heartmarshall/my-english/internal/service/word"
)

// RepositoryFactory создаёт репозитории с указанным Querier.
type RepositoryFactory struct{}

// NewRepositoryFactory создаёт новую фабрику репозиториев.
func NewRepositoryFactory() *RepositoryFactory {
	return &RepositoryFactory{}
}

// Words создаёт репозиторий слов.
func (f *RepositoryFactory) Words(q database.Querier) wordservice.WordRepository {
	return word.New(q)
}

// Meanings создаёт репозиторий значений.
func (f *RepositoryFactory) Meanings(q database.Querier) wordservice.MeaningRepository {
	return meaning.New(q, meaning.WithClock(database.RealClock{}))
}

// Examples создаёт репозиторий примеров.
func (f *RepositoryFactory) Examples(q database.Querier) wordservice.ExampleRepository {
	return example.New(q)
}

// Tags создаёт репозиторий тегов.
func (f *RepositoryFactory) Tags(q database.Querier) wordservice.TagRepository {
	return tag.New(q)
}

// MeaningTags создаёт репозиторий связей meaning-tag.
func (f *RepositoryFactory) MeaningTags(q database.Querier) wordservice.MeaningTagRepository {
	return meaningtag.New(q)
}

// Translations создаёт репозиторий переводов.
func (f *RepositoryFactory) Translations(q database.Querier) wordservice.TranslationRepository {
	return translation.New(q)
}

// Compile-time check
var _ wordservice.RepositoryFactory = (*RepositoryFactory)(nil)
