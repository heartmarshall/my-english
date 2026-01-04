package study

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/model"
	"github.com/heartmarshall/my-english/internal/service"
	"github.com/heartmarshall/my-english/internal/service/study/srs"
)

type ReviewResult struct {
	Card          model.Card
	SRS           model.SRSState
	StatusChanged bool
}

// Review обрабатывает ответ пользователя на карточку.
func (s *Service) Review(ctx context.Context, cardID uuid.UUID, grade int, durationMs int) (*ReviewResult, error) {
	if grade < 1 || grade > 5 {
		return nil, service.ErrInvalidGrade
	}

	var result *ReviewResult

	err := s.txManager.RunInTx(ctx, func(ctx context.Context, tx database.Querier) error {
		srsRepo := s.repos.SRS(tx)
		logRepo := s.repos.Review(tx)
		cardRepo := s.repos.Card(tx)

		// 1. Получаем текущее состояние SRS с блокировкой (FOR UPDATE)
		// Это предотвращает race conditions при параллельных review запросах
		currentState, err := srsRepo.GetByCardIDForUpdate(ctx, cardID)
		if err != nil && !database.IsNotFoundError(err) {
			return err
		}
		if currentState == nil {
			currentState = &model.SRSState{
				CardID:        cardID,
				Status:        model.LearningStatusNew,
				AlgorithmData: map[string]any{},
			}
		}

		// 2. Рассчитываем новое состояние
		now := time.Now()
		input := srs.Input{
			Status:        currentState.Status,
			AlgorithmData: currentState.AlgorithmData,
			Grade:         grade,
			Now:           now,
		}

		output := s.algo.Calculate(input)

		// 3. Обновляем SRS State
		newState := &model.SRSState{
			CardID:        cardID,
			Status:        output.Status,
			DueDate:       &output.NextReviewAt,
			AlgorithmData: output.AlgorithmData,
			LastReviewAt:  &now,
		}

		updatedState, err := srsRepo.Upsert(ctx, newState)
		if err != nil {
			return fmt.Errorf("failed to update srs state: %w", err)
		}

		// 4. Пишем лог (историю)
		// Используем указатель для опционального duration
		var dur *int
		if durationMs > 0 {
			dur = &durationMs
		}

		err = logRepo.Create(ctx, &model.ReviewLog{
			CardID:      cardID,
			Grade:       grade,
			DurationMs:  dur,
			ReviewedAt:  now,
			StateBefore: currentState.AlgorithmData,
			StateAfter:  output.AlgorithmData,
		})
		if err != nil {
			return fmt.Errorf("failed to create review log: %w", err)
		}

		// 5. Загружаем саму карточку для ответа
		card, err := cardRepo.GetByID(ctx, cardID)
		if err != nil {
			return err
		}

		result = &ReviewResult{
			Card:          *card,
			SRS:           *updatedState,
			StatusChanged: currentState.Status != output.Status,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
