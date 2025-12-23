// Package service содержит бизнес-логику приложения (use cases).
// Сервисы принимают и возвращают доменные модели из пакета model.
// На сервисный слой опирается транспортный слой (GraphQL resolvers).
package service

import "errors"

// Бизнес-ошибки сервисного слоя.
var (
	// ErrWordNotFound возвращается, когда слово не найдено.
	ErrWordNotFound = errors.New("word not found")

	// ErrWordAlreadyExists возвращается при попытке создать дубликат слова.
	ErrWordAlreadyExists = errors.New("word already exists")

	// ErrMeaningNotFound возвращается, когда значение не найдено.
	ErrMeaningNotFound = errors.New("meaning not found")

	// ErrInvalidInput возвращается при невалидных входных данных.
	ErrInvalidInput = errors.New("invalid input")

	// ErrInvalidGrade возвращается при невалидной оценке (должна быть 1-5).
	ErrInvalidGrade = errors.New("grade must be between 1 and 5")
)
