package study

// Service содержит бизнес-логику для изучения слов (SRS).
type Service struct {
	meanings MeaningRepository
	srs      MeaningSRSRepository
	clock    Clock
}

// Deps — зависимости для создания сервиса.
type Deps struct {
	Meanings MeaningRepository
	SRS      MeaningSRSRepository
	Clock    Clock // опционально, по умолчанию RealClock
}

// New создаёт новый сервис.
func New(deps Deps) *Service {
	clock := deps.Clock
	if clock == nil {
		clock = RealClock{}
	}

	return &Service{
		meanings: deps.Meanings,
		srs:      deps.SRS,
		clock:    clock,
	}
}
