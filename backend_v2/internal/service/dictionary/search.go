package dictionary

import (
	"context"
	"log/slog"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/heartmarshall/my-english/internal/model"
	ctxlog "github.com/heartmarshall/my-english/pkg/context"
	"golang.org/x/sync/errgroup"
)

// Search ищет слова по запросу query.
// Логика работы (Read-Through Cache):
// 1. Ищет в локальной БД (нечеткий поиск по триграммам).
// 2. Если точного совпадения нет — параллельно опрашивает всех внешних провайдеров.
// 3. Найденные вовне данные сохраняет в БД.
// 4. Возвращает объединенный результат.
func (s *Service) Search(ctx context.Context, query string) ([]model.Lexeme, error) {
	query = strings.TrimSpace(strings.ToLower(query))
	if query == "" {
		return nil, nil
	}

	logger := ctxlog.L(ctx)

	// 1. Ищем локально (используем пул соединений вне транзакции)
	// Лимит 10, чтобы не забивать выдачу похожими словами, если их слишком много
	localResults, err := s.repos.Lexeme(s.txManager.Q()).SearchFuzzy(ctx, query, 10)
	if err != nil {
		return nil, err
	}

	// Проверяем, есть ли среди локальных результатов точное совпадение
	hasExactMatch := false
	for _, l := range localResults {
		if l.TextNormalized == query {
			hasExactMatch = true
			break
		}
	}

	// 2. Если точного совпадения нет, запускаем внешние источники
	// Если провайдеров нет, просто возвращаем локальные результаты
	if !hasExactMatch && len(s.providers) > 0 {
		var mu sync.Mutex
		foundLexemes := make([]*model.Lexeme, 0)

		// Создаем группу горутин с контекстом
		g, groupCtx := errgroup.WithContext(ctx)

		for _, provider := range s.providers {
			// Замыкание переменных для использования внутри горутины
			p := provider

			g.Go(func() error {
				// Используем groupCtx, чтобы отменить остальные запросы, если кто-то вернет критическую ошибку
				// (хотя мы здесь игнорируем ошибки провайдеров, чтобы не ломать общий поиск)
				imported, err := p.Fetch(groupCtx, query)
				if err != nil {
					// Логируем ошибку, но возвращаем nil, чтобы errgroup не отменял остальные горутины
					logger.Warn("dictionary provider failed",
						slog.String("provider", p.SourceSlug()),
						slog.String("query", query),
						slog.String("error", err.Error()))
					return nil
				}

				if imported == nil {
					return nil // Провайдер ничего не нашел
				}

				// Сохраняем результат в БД.
				// SaveImportedWord открывает свою независимую транзакцию, это безопасно в конкурентной среде.
				lexeme, err := s.SaveImportedWord(groupCtx, imported, p.SourceSlug())
				if err != nil {
					logger.Error("failed to save imported word",
						slog.String("provider", p.SourceSlug()),
						slog.String("error", err.Error()))
					return nil
				}

				// Добавляем результат в общий список (потокобезопасно)
				mu.Lock()
				foundLexemes = append(foundLexemes, lexeme)
				mu.Unlock()
				return nil
			})
		}

		// Ждем завершения всех провайдеров
		if err := g.Wait(); err != nil {
			// В текущей логике сюда мы попадем только при критических системных ошибках,
			// так как ошибки провайдеров мы подавили.
			return nil, err
		}

		// 3. Объединяем результаты
		if len(foundLexemes) > 0 {
			// Создаем мапу существующих ID, чтобы не добавить дубликаты
			// (если локальный поиск вернул нечеткое, а API вернул точное, которое совпало с нечетким)
			existingIDs := make(map[uuid.UUID]bool)
			for _, l := range localResults {
				existingIDs[l.ID] = true
			}

			// Добавляем новые найденные слова в НАЧАЛО списка
			newItems := make([]model.Lexeme, 0, len(foundLexemes))
			for _, ptr := range foundLexemes {
				if !existingIDs[ptr.ID] {
					newItems = append(newItems, *ptr)
				}
			}

			localResults = append(newItems, localResults...)
		}
	}

	return localResults, nil
}
