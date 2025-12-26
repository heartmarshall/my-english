package app

import (
	"context"

	"github.com/heartmarshall/my-english/graph"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/dictionary"
	"github.com/heartmarshall/my-english/internal/database/example"
	inboxrepo "github.com/heartmarshall/my-english/internal/database/inbox"
	"github.com/heartmarshall/my-english/internal/database/meaning"
	"github.com/heartmarshall/my-english/internal/database/meaningtag"
	"github.com/heartmarshall/my-english/internal/database/tag"
	"github.com/heartmarshall/my-english/internal/database/translation"
	"github.com/heartmarshall/my-english/internal/database/word"
	inboxservice "github.com/heartmarshall/my-english/internal/service/inbox"
	"github.com/heartmarshall/my-english/internal/service/loader"
	"github.com/heartmarshall/my-english/internal/service/study"
	wordservice "github.com/heartmarshall/my-english/internal/service/word"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TxRunner реализует интерфейс wordservice.TxRunner.
type TxRunner struct {
	txManager *database.TxManager
}

// NewTxRunner создаёт новый TxRunner.
func NewTxRunner(pool *pgxpool.Pool) *TxRunner {
	return &TxRunner{
		txManager: database.NewTxManager(pool),
	}
}

// RunInTx выполняет функцию в транзакции.
func (r *TxRunner) RunInTx(ctx context.Context, fn func(ctx context.Context, tx database.Querier) error) error {
	return r.txManager.RunInTx(ctx, fn)
}

// Compile-time check
var _ wordservice.TxRunner = (*TxRunner)(nil)

// Repositories содержит все репозитории.
type Repositories struct {
	Words        *word.Repo
	Meanings     *meaning.Repo
	Examples     *example.Repo
	Tags         *tag.Repo
	MeaningTag   *meaningtag.Repo
	Translations *translation.Repo
	Dictionary   *dictionary.Repo
	Inbox        *inboxrepo.Repo
}

// Services содержит все сервисы.
type Services struct {
	Words  *wordservice.Service
	Study  *study.Service
	Loader *loader.Service
	Inbox  *inboxservice.Service
}

// Dependencies содержит все зависимости приложения.
type Dependencies struct {
	DB           *pgxpool.Pool
	Repositories *Repositories
	Services     *Services
	Resolver     *graph.Resolver
}

// NewDependencies создаёт все зависимости приложения.
func NewDependencies(pool *pgxpool.Pool) *Dependencies {
	// Репозитории
	repos := newRepositories(pool)

	// TxRunner и RepositoryFactory
	txRunner := NewTxRunner(pool)
	repoFactory := NewRepositoryFactory()

	// Сервисы
	services := newServices(repos, txRunner, repoFactory)

	// GraphQL Resolver (использует только сервисы)
	resolver := newResolver(services)

	return &Dependencies{
		DB:           pool,
		Repositories: repos,
		Services:     services,
		Resolver:     resolver,
	}
}

func newRepositories(pool *pgxpool.Pool) *Repositories {
	return &Repositories{
		Words:        word.New(pool),
		Meanings:     meaning.New(pool, meaning.WithClock(database.RealClock{})),
		Examples:     example.New(pool),
		Tags:         tag.New(pool),
		MeaningTag:   meaningtag.New(pool),
		Translations: translation.New(pool),
		Dictionary:   dictionary.New(pool, dictionary.WithClock(database.RealClock{})),
		Inbox:        inboxrepo.New(pool, inboxrepo.WithClock(database.RealClock{})),
	}
}

func newServices(repos *Repositories, txRunner *TxRunner, repoFactory *RepositoryFactory) *Services {
	// Word Service
	wordSvc := wordservice.New(wordservice.Deps{
		Words:        repos.Words,
		Meanings:     repos.Meanings,
		Examples:     repos.Examples,
		Tags:         repos.Tags,
		MeaningTag:   repos.MeaningTag,
		Translations: repos.Translations,
		Dictionary:   repos.Dictionary,
		TxRunner:     txRunner,
		RepoFactory:  repoFactory,
	})

	// Study Service — используем SRS адаптер
	srsAdapter := NewSRSAdapter(repos.Meanings)
	studySvc := study.New(study.Deps{
		Meanings: repos.Meanings,
		SRS:      srsAdapter,
		Clock:    study.RealClock{},
	})

	// Loader Service для DataLoaders
	loaderSvc := loader.New(loader.Deps{
		Meanings:     repos.Meanings,
		Examples:     repos.Examples,
		Tags:         repos.Tags,
		MeaningTags:  repos.MeaningTag,
		Translations: repos.Translations,
	})

	// Inbox Service
	inboxSvc := inboxservice.New(inboxservice.Deps{
		Inbox: repos.Inbox,
	})

	return &Services{
		Words:  wordSvc,
		Study:  studySvc,
		Loader: loaderSvc,
		Inbox:  inboxSvc,
	}
}

func newResolver(services *Services) *graph.Resolver {
	// Resolver использует только сервисы
	return graph.NewResolver(graph.Deps{
		Words: services.Words,
		Study: services.Study,
		Inbox: services.Inbox,
	})
}
