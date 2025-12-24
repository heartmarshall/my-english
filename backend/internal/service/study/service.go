package study

import (
	"github.com/heartmarshall/my-english/internal/service/study/srs"
)

// Service содержит бизнес-логику для изучения слов (SRS).
type Service struct {
	meanings MeaningRepository
	srs      MeaningSRSRepository
	clock    Clock
	strategy srs.Strategy
}

// Deps — зависимости для создания сервиса.
type Deps struct {
	Meanings MeaningRepository
	SRS      MeaningSRSRepository
	Clock    Clock        // опционально, по умолчанию RealClock
	Strategy srs.Strategy // опционально, по умолчанию SM2Strategy
}

// New создаёт новый сервис.
func New(deps Deps) *Service {
	clock := deps.Clock
	if clock == nil {
		clock = RealClock{}
	}

	strategy := deps.Strategy
	if strategy == nil {
		strategy = srs.NewSM2Strategy()
	}

	return &Service{
		meanings: deps.Meanings,
		srs:      deps.SRS,
		clock:    clock,
		strategy: strategy,
	}
}
