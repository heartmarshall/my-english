package srs

import (
	"time"

	"github.com/heartmarshall/my-english/internal/model"
)

// SM-2 Algorithm constants
// https://en.wikipedia.org/wiki/SuperMemo#Description_of_SM-2_algorithm
const (
	minEaseFactor     = 1.3
	defaultEaseFactor = 2.5
)

// SM2Strategy реализует алгоритм SM-2 для интервального повторения.
type SM2Strategy struct{}

// NewSM2Strategy создаёт новую стратегию SM-2.
func NewSM2Strategy() *SM2Strategy {
	return &SM2Strategy{}
}

// Calculate вычисляет новые SRS параметры на основе алгоритма SM-2.
// grade: 1-5, где 1 = не помню, 5 = отлично помню
//
// Использует упрощённый алгоритм SM-2:
// - grade < 3: сбрасываем интервал, уменьшаем ease factor
// - grade >= 3: увеличиваем интервал, корректируем ease factor
func (s *SM2Strategy) Calculate(meaning *model.Meaning, grade int, now time.Time) *Update {
	// Текущие значения или дефолты
	currentInterval := 1
	if meaning.Interval != nil {
		currentInterval = *meaning.Interval
	}

	currentEase := defaultEaseFactor
	if meaning.EaseFactor != nil {
		currentEase = *meaning.EaseFactor
	}

	currentReviewCount := 0
	if meaning.ReviewCount != nil {
		currentReviewCount = *meaning.ReviewCount
	}

	var (
		newInterval int
		newEase     float64
		newStatus   model.LearningStatus
	)

	if grade < 3 {
		// Неудачный ответ — сбрасываем прогресс
		newInterval = 1
		newEase = max(minEaseFactor, currentEase-0.2)
		newStatus = model.LearningStatusLearning
	} else {
		// Успешный ответ — увеличиваем интервал
		if meaning.LearningStatus == model.LearningStatusNew {
			// Первое повторение нового слова
			newInterval = 1
		} else if currentInterval == 1 {
			newInterval = 6
		} else {
			newInterval = int(float64(currentInterval) * currentEase)
		}

		// Корректируем ease factor по формуле SM-2
		newEase = currentEase + (0.1 - float64(5-grade)*(0.08+float64(5-grade)*0.02))
		newEase = max(minEaseFactor, newEase)

		// Определяем статус
		if newInterval >= 21 {
			newStatus = model.LearningStatusMastered
		} else if meaning.LearningStatus == model.LearningStatusNew || newInterval < 6 {
			// Первое повторение нового слова или маленький интервал -> LEARNING
			newStatus = model.LearningStatusLearning
		} else {
			// Интервал >= 6 и < 21 -> REVIEW
			newStatus = model.LearningStatusReview
		}
	}

	nextReview := now.Add(time.Duration(newInterval) * 24 * time.Hour)
	newReviewCount := currentReviewCount + 1

	return &Update{
		LearningStatus: newStatus,
		NextReviewAt:   &nextReview,
		Interval:       &newInterval,
		EaseFactor:     &newEase,
		ReviewCount:    &newReviewCount,
	}
}
