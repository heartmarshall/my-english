package srs

import (
	"math"
	"time"

	"github.com/heartmarshall/my-english/internal/model"
)

// Constants for SM-2
const (
	defaultEaseFactor = 2.5
	minEaseFactor     = 1.3
	defaultInterval   = 1 // 1 день для новых
)

// Input — входные данные для алгоритма.
type Input struct {
	Status        model.LearningStatus
	AlgorithmData map[string]any // Сырые данные из БД (interval, ease_factor)
	Grade         int            // Оценка 1-5
	Now           time.Time
}

// Output — результат работы алгоритма.
type Output struct {
	Status        model.LearningStatus
	NextReviewAt  time.Time
	AlgorithmData map[string]any // Обновленные данные для БД
}

// Algorithm — интерфейс стратегии интервального повторения.
type Algorithm interface {
	Calculate(input Input) Output
}

// SM2 — реализация алгоритма SuperMemo-2.
type SM2 struct{}

func NewSM2() *SM2 {
	return &SM2{}
}

func (a *SM2) Calculate(in Input) Output {
	// 1. Извлекаем текущее состояние из JSONB
	interval := getInt(in.AlgorithmData, "interval", 0)
	easeFactor := getFloat(in.AlgorithmData, "ease_factor", defaultEaseFactor)
	reviewCount := getInt(in.AlgorithmData, "review_count", 0)

	var nextInterval int
	var nextStatus model.LearningStatus

	// 2. Логика SM-2
	if in.Grade < 3 {
		// Забыл: сброс
		reviewCount = 0
		nextInterval = 1
		nextStatus = model.LearningStatusLearning
		// Ease factor не меняется или немного уменьшается (опционально)
	} else {
		// Вспомнил
		reviewCount++

		if interval == 0 {
			nextInterval = 1
		} else if interval == 1 {
			nextInterval = 6
		} else {
			nextInterval = int(math.Ceil(float64(interval) * easeFactor))
		}

		// Корректировка Ease Factor
		// EF' = EF + (0.1 - (5-q) * (0.08 + (5-q) * 0.02))
		delta := 5 - float64(in.Grade)
		easeFactor = easeFactor + (0.1 - delta*(0.08+delta*0.02))
		if easeFactor < minEaseFactor {
			easeFactor = minEaseFactor
		}

		// Переход статусов
		if nextInterval > 21 {
			nextStatus = model.LearningStatusMastered
		} else {
			nextStatus = model.LearningStatusReview
		}
	}

	// 3. Формируем результат
	return Output{
		Status:       nextStatus,
		NextReviewAt: in.Now.AddDate(0, 0, nextInterval),
		AlgorithmData: map[string]any{
			"interval":     nextInterval,
			"ease_factor":  easeFactor,
			"review_count": reviewCount,
		},
	}
}

// Helpers для безопасного извлечения из map[string]any
func getInt(m map[string]any, key string, def int) int {
	if val, ok := m[key]; ok {
		// JSON числа часто приходят как float64
		if f, ok := val.(float64); ok {
			return int(f)
		}
		if i, ok := val.(int); ok {
			return i
		}
	}
	return def
}

func getFloat(m map[string]any, key string, def float64) float64 {
	if val, ok := m[key]; ok {
		if f, ok := val.(float64); ok {
			return f
		}
	}
	return def
}
