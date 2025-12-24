package study

import (
	"context"
	"errors"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/model"
	"github.com/heartmarshall/my-english/internal/service"
)

// ReviewMeaning обрабатывает оценку пользователя и обновляет SRS данные.
// grade: 1-5, где 1 = не помню, 5 = отлично помню
//
// Использует стратегию SRS для расчета новых параметров интервального повторения.
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

	// Рассчитываем новые SRS параметры используя стратегию
	now := s.clock.Now()
	srsUpdate := s.strategy.Calculate(&meaning, grade, now)

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
