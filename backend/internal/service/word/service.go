package word

// Service содержит бизнес-логику для работы со словами.
type Service struct {
	words        WordRepository
	meanings     MeaningRepository
	examples     ExampleRepository
	tags         TagRepository
	meaningTag   MeaningTagRepository
	translations TranslationRepository

	dictionary DictionaryRepository
	fetcher    DictionaryFetcher // <-- Внешний источник данных

	txRunner    TxRunner
	repoFactory RepositoryFactory
}

// Deps — зависимости для создания сервиса.
type Deps struct {
	Words        WordRepository
	Meanings     MeaningRepository
	Examples     ExampleRepository
	Tags         TagRepository
	MeaningTag   MeaningTagRepository
	Translations TranslationRepository
	Dictionary   DictionaryRepository

	Fetcher DictionaryFetcher // <-- Инжектим реализацию клиента сюда

	TxRunner    TxRunner
	RepoFactory RepositoryFactory
}

// New создаёт новый сервис.
func New(deps Deps) *Service {
	return &Service{
		words:        deps.Words,
		meanings:     deps.Meanings,
		examples:     deps.Examples,
		tags:         deps.Tags,
		meaningTag:   deps.MeaningTag,
		translations: deps.Translations,
		dictionary:   deps.Dictionary,
		fetcher:      deps.Fetcher, // <-- Сохраняем
		txRunner:     deps.TxRunner,
		repoFactory:  deps.RepoFactory,
	}
}
