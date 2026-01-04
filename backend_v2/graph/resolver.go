package graph

import (
	"github.com/heartmarshall/my-english/internal/service/card"
	"github.com/heartmarshall/my-english/internal/service/dictionary"
	"github.com/heartmarshall/my-english/internal/service/inbox"
	"github.com/heartmarshall/my-english/internal/service/loader"
	"github.com/heartmarshall/my-english/internal/service/stats"
	"github.com/heartmarshall/my-english/internal/service/study"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	CardService       *card.Service
	DictionaryService *dictionary.Service
	InboxService      *inbox.Service
	LoaderService     *loader.Service
	StatsService      *stats.Service
	StudyService      *study.Service
}

// Card returns CardResolver implementation.
func (r *Resolver) Card() CardResolver { return &cardResolver{r} }

// Lexeme returns LexemeResolver implementation.
func (r *Resolver) Lexeme() LexemeResolver { return &lexemeResolver{r} }

// Sense returns SenseResolver implementation.
func (r *Resolver) Sense() SenseResolver { return &senseResolver{r} }
