package study

import (
	"context"
	"errors"
	"time"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/model"
	"github.com/heartmarshall/my-english/internal/service"
)

// SM-2 Algorithm constants
// https://en.wikipedia.org/wiki/SuperMemo#Description_of_SM-2_algorithm
const (
	minEaseFactor     = 1.3
	defaultEaseFactor = 2.5
)

// ReviewMeaning обрабатывает оценку пользователя и обновляет SRS данные.
// grade: 1-5, где 1 = не помню, 5 = отлично помню
//
// Использует упрощённый алгоритм SM-2:
// - grade < 3: сбрасываем интервал, уменьшаем ease factor
// - grade >= 3: увеличиваем интервал, корректируем ease factor
func (s *Service) ReviewMeaning(ctx context.Context, meaningID int64, grade int) (model.Meaning, error) {
	// Валидация
	if grade < 1 || grade > 5 {
		return model.Meaning{}, service.ErrInvalidGrade
	}

	// Получаем текущее состояние
	meaning, err := s.meanings.GetByID(ctx, meaningID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return model.Meaning{}, service.ErrMeaningNotFound
		}
		return model.Meaning{}, err
	}

	// Рассчитываем новые SRS параметры
	srsUpdate := s.calculateSRS(&meaning, grade)

	// Обновляем в БД
	if err := s.srs.UpdateSRS(ctx, meaningID, srsUpdate); err != nil {
		return model.Meaning{}, err
	}

	// Обновляем модель для возврата
	meaning.LearningStatus = srsUpdate.LearningStatus
	meaning.NextReviewAt = srsUpdate.NextReviewAt
	meaning.Interval = srsUpdate.Interval
	meaning.EaseFactor = srsUpdate.EaseFactor
	meaning.ReviewCount = srsUpdate.ReviewCount

	return meaning, nil
}

// calculateSRS рассчитывает новые SRS параметры на основе оценки.
func (s *Service) calculateSRS(meaning *model.Meaning, grade int) *SRSUpdate {
	now := s.clock.Now()

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
		} else {
			newStatus = model.LearningStatusReview
		}
	}

	nextReview := now.Add(time.Duration(newInterval) * 24 * time.Hour)
	newReviewCount := currentReviewCount + 1

	return &SRSUpdate{
		LearningStatus: newStatus,
		NextReviewAt:   &nextReview,
		Interval:       &newInterval,
		EaseFactor:     &newEase,
		ReviewCount:    &newReviewCount,
	}
}
