package card

import (
	"github.com/google/uuid"
	"github.com/heartmarshall/my-english/internal/model"
)

type CreateCardInput struct {
	// Вариант А: Ссылка на словарь
	SenseID *uuid.UUID

	// Вариант Б: Полностью свое слово
	CustomText *string

	// Общие поля / Переопределения
	CustomTranscription *string
	CustomTranslations  []string
	CustomNote          *string
	CustomImageURL      *string

	// Теги (строки, например ["IT", "Verbs"])
	Tags []string
}

type UpdateCardInput struct {
	CustomTranscription *string
	CustomTranslations  []string
	CustomNote          *string
	CustomImageURL      *string
	Tags                []string
}

// Filter используется для поиска карточек пользователя
type Filter struct {
	Search   *string
	Tags     []string
	Statuses []model.LearningStatus
}
