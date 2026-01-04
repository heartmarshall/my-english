// Package service содержит бизнес-логику приложения.
package service

import "errors"

// Ошибки сервисного слоя.
var (
	// ErrWordNotFound — слово не найдено.
	ErrWordNotFound = errors.New("word not found")

	// ErrMeaningNotFound — значение не найдено.
	ErrMeaningNotFound = errors.New("meaning not found")

	// ErrWordAlreadyExists — слово уже существует.
	ErrWordAlreadyExists = errors.New("word already exists")

	// ErrInvalidInput — невалидные входные данные.
	ErrInvalidInput = errors.New("invalid input")

	// ErrInvalidGrade — невалидная оценка (должна быть 1-5).
	ErrInvalidGrade = errors.New("invalid grade")

	// ErrCardNotFound — карточка не найдена.
	ErrCardNotFound = errors.New("card not found")

	// ErrLexemeNotFound — лексема не найдена.
	ErrLexemeNotFound = errors.New("lexeme not found")

	// ErrSenseNotFound — смысл не найден.
	ErrSenseNotFound = errors.New("sense not found")
)

