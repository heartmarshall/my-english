package srs

import (
	"time"

	"github.com/heartmarshall/my-english/internal/model"
)

// Update содержит поля для обновления SRS данных.
type Update struct {
	LearningStatus model.LearningStatus
	NextReviewAt   *time.Time
	Interval       *int
	EaseFactor     *float64
	ReviewCount    *int
}

// Strategy определяет интерфейс для алгоритмов интервального повторения.
// Различные реализации могут использовать SM-2, FSRS, Anki-style и т.д.
type Strategy interface {
	// Calculate вычисляет новые SRS параметры на основе текущего состояния и оценки пользователя.
	// grade: 1-5, где 1 = не помню, 5 = отлично помню
	Calculate(meaning *model.Meaning, grade int, now time.Time) *Update
}
