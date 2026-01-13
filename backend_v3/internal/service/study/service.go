package study

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/repository"
	"github.com/heartmarshall/my-english/internal/model"
)

type Service struct {
	repos *repository.Registry
	tx    *database.TxManager
}

func NewService(repos *repository.Registry, tx *database.TxManager) *Service {
	return &Service{
		repos: repos,
		tx:    tx,
	}
}

// ReviewCardInput данные для ответа на карточку.
type ReviewCardInput struct {
	CardID     uuid.UUID
	Grade      model.ReviewGrade
	DurationMs *int
	ReviewedAt time.Time // Обычно Now, но полезно для тестов или оффлайн-синхронизации
}

// ReviewResult возвращается после успешного повторения.
type ReviewResult struct {
	Card         model.Card
	ReviewLog    model.ReviewLog
	NextReviewAt time.Time
}

// ReviewCard обрабатывает ответ пользователя на карточку.
// Это критическая секция, требующая блокировки строки.
func (s *Service) ReviewCard(ctx context.Context, input ReviewCardInput) (*ReviewResult, error) {
	if input.ReviewedAt.IsZero() {
		input.ReviewedAt = time.Now()
	}

	var result ReviewResult

	err := s.tx.RunInTx(ctx, func(ctx context.Context, _ database.Querier) error {
		// 1. Получаем карточку с блокировкой FOR UPDATE
		// Это предотвращает race conditions, если пользователь дважды нажмет кнопку,
		// или если параллельно придет запрос на обновление.
		card, err := s.repos.Cards.GetByIDForUpdate(ctx, input.CardID)
		if err != nil {
			return err
		}

		// 2. Рассчитываем новые параметры SRS (Чистая функция)
		srsCalc := calculateSRS(card, input.Grade, input.ReviewedAt)

		// 3. Обновляем карточку
		// Используем UpdateSRSFields для оптимизации (обновляем только нужные поля)
		err = s.repos.Cards.UpdateSRSFields(
			ctx,
			card.ID,
			srsCalc.Status,
			&srsCalc.NextReviewAt,
			srsCalc.IntervalDays,
			srsCalc.EaseFactor,
		)
		if err != nil {
			return err
		}

		// 4. Записываем лог повторения (Review Log)
		logEntry := &model.ReviewLog{
			CardID:     card.ID,
			Grade:      input.Grade,
			DurationMs: input.DurationMs,
			ReviewedAt: input.ReviewedAt,
		}
		createdLog, err := s.repos.ReviewLogs.Create(ctx, logEntry)
		if err != nil {
			return err
		}

		// Подготавливаем результат (обновляем поля в объекте card для возврата)
		card.Status = srsCalc.Status
		card.NextReviewAt = &srsCalc.NextReviewAt
		card.IntervalDays = srsCalc.IntervalDays
		card.EaseFactor = srsCalc.EaseFactor
		card.UpdatedAt = time.Now() // Примерно

		result = ReviewResult{
			Card:         *card,
			ReviewLog:    *createdLog,
			NextReviewAt: srsCalc.NextReviewAt,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetStudyQueue возвращает очередь карточек для изучения.
// Логика выборки инкапсулирована в репозитории, сервис просто проксирует.
func (s *Service) GetStudyQueue(ctx context.Context, limit int) ([]model.DictionaryEntry, error) {
	// 1. Получаем карточки, которые пора повторять
	cards, err := s.repos.Cards.GetDueCards(ctx, time.Now(), limit)
	if err != nil {
		return nil, err
	}

	if len(cards) == 0 {
		return []model.DictionaryEntry{}, nil
	}

	// 2. Собираем ID слов
	entryIDs := make([]any, len(cards))
	// Map для быстрого присоединения карточки к слову (если нужно будет возвращать DTO)
	// Но здесь мы возвращаем DictionaryEntry.
	// В GraphQL резолверах мы подтянем Card через Dataloader,
	// но для Study режима удобно сразу вернуть Entry.

	// ВАЖНО: Архитектурный момент.
	// Метод возвращает DictionaryEntry. Но фронту нужна и Card.
	// GraphQL позволяет запросить `studyQueue { card { ... } }`.
	// Поэтому здесь достаточно вернуть Entry.

	for i, c := range cards {
		entryIDs[i] = c.EntryID
	}

	// 3. Загружаем слова
	// Порядок в IN (...) не гарантирован, поэтому нужно отсортировать в соответствии с очередью cards
	entries, err := s.repos.Dictionary.ListByIDs(ctx, "id", entryIDs)
	if err != nil {
		return nil, err
	}

	// Восстанавливаем порядок
	entriesMap := make(map[uuid.UUID]model.DictionaryEntry, len(entries))
	for _, e := range entries {
		entriesMap[e.ID] = e
	}

	orderedEntries := make([]model.DictionaryEntry, 0, len(cards))
	for _, c := range cards {
		if e, ok := entriesMap[c.EntryID]; ok {
			// Небольшой хак: прокидываем карточку внутрь Entry, если модель это позволяет,
			// или полагаемся на Dataloader.
			// В model.DictionaryEntry есть поле `card *Card` (см. анализ schema.graphqls / models).
			// В Go-модели DictionaryEntry поля Card нет (она чистая DB модель).
			// Значит, связка произойдет на уровне GraphQL Resolver'а.
			orderedEntries = append(orderedEntries, e)
		}
	}

	return orderedEntries, nil
}

func (s *Service) GetDashboardStats(ctx context.Context) (*repository.DashboardStatsDTO, error) {
	return s.repos.Cards.GetDashboardStats(ctx)
}
