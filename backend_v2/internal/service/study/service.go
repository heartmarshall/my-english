package study

import (
	"github.com/heartmarshall/my-english/internal/database"
	factory "github.com/heartmarshall/my-english/internal/database/repository/factory"
	"github.com/heartmarshall/my-english/internal/service/study/srs"
)

type Service struct {
	repos     *factory.Factory
	txManager *database.TxManager
	algo      srs.Algorithm
}

type Deps struct {
	Repos     *factory.Factory
	TxManager *database.TxManager
	Algorithm srs.Algorithm // Опционально, иначе SM2
}

func New(deps Deps) *Service {
	algo := deps.Algorithm
	if algo == nil {
		algo = srs.NewSM2()
	}

	return &Service{
		repos:     deps.Repos,
		txManager: deps.TxManager,
		algo:      algo,
	}
}
