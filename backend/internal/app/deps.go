package app

import (
	"database/sql"

	"github.com/heartmarshall/my-english/graph"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/example"
	"github.com/heartmarshall/my-english/internal/database/meaning"
	"github.com/heartmarshall/my-english/internal/database/meaningtag"
	"github.com/heartmarshall/my-english/internal/database/tag"
	"github.com/heartmarshall/my-english/internal/database/word"
	"github.com/heartmarshall/my-english/internal/service/study"
	wordservice "github.com/heartmarshall/my-english/internal/service/word"
)

// Repositories содержит все репозитории.
type Repositories struct {
	Words      *word.Repo
	Meanings   *meaning.Repo
	Examples   *example.Repo
	Tags       *tag.Repo
	MeaningTag *meaningtag.Repo
}

// Services содержит все сервисы.
type Services struct {
	Words *wordservice.Service
	Study *study.Service
}

// Dependencies содержит все зависимости приложения.
type Dependencies struct {
	DB           *sql.DB
	Repositories *Repositories
	Services     *Services
	Resolver     *graph.Resolver
}

// NewDependencies создаёт все зависимости приложения.
func NewDependencies(db *sql.DB) *Dependencies {
	// Репозитории
	repos := newRepositories(db)

	// Сервисы
	services := newServices(repos)

	// GraphQL Resolver
	resolver := newResolver(services, repos)

	return &Dependencies{
		DB:           db,
		Repositories: repos,
		Services:     services,
		Resolver:     resolver,
	}
}

func newRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Words:      word.New(db),
		Meanings:   meaning.New(db, meaning.WithClock(database.RealClock{})),
		Examples:   example.New(db),
		Tags:       tag.New(db),
		MeaningTag: meaningtag.New(db),
	}
}

func newServices(repos *Repositories) *Services {
	// Word Service
	wordSvc := wordservice.New(wordservice.Deps{
		Words:      repos.Words,
		Meanings:   repos.Meanings,
		Examples:   repos.Examples,
		Tags:       repos.Tags,
		MeaningTag: repos.MeaningTag,
	})

	// Study Service — используем SRS адаптер
	srsAdapter := NewSRSAdapter(repos.Meanings)
	studySvc := study.New(study.Deps{
		Meanings: repos.Meanings,
		SRS:      srsAdapter,
		Clock:    study.RealClock{},
	})

	return &Services{
		Words: wordSvc,
		Study: studySvc,
	}
}

func newResolver(services *Services, repos *Repositories) *graph.Resolver {
	// Используем TagLoader адаптер
	tagLoader := NewTagLoaderAdapter(repos.Tags, repos.MeaningTag)

	return graph.NewResolver(graph.Deps{
		Words:    services.Words,
		Study:    services.Study,
		Examples: repos.Examples,
		Tags:     tagLoader,
	})
}
