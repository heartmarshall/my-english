package suggestion

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	ctx_pkg "github.com/heartmarshall/my-english/pkg/context"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/singleflight"
)

type Service struct {
	providers map[string]Provider
	sf        singleflight.Group
}

// NewService создает сервис и регистрирует провайдеров.
func NewService(providers ...Provider) *Service {
	pMap := make(map[string]Provider, len(providers))
	for _, p := range providers {
		pMap[p.Slug()] = p
	}
	return &Service{
		providers: pMap,
	}
}

// FetchSuggestions получает подсказки из указанных источников.
// Реализует паттерн Scatter-Gather.
func (s *Service) FetchSuggestions(ctx context.Context, text string, sources []string) ([]Result, error) {
	// 1. Дедупликация запросов (Singleflight)
	// Ключ зависит от текста и списка источников.
	key := fmt.Sprintf("%s:%v", text, sources)

	// Singleflight возвращает (interface{}, error, shared bool).
	val, err, _ := s.sf.Do(key, func() (interface{}, error) {
		return s.fetchInternal(ctx, text, sources)
	})

	if err != nil {
		return nil, err
	}

	return val.([]Result), nil
}

// fetchInternal выполняет реальную логику запросов.
func (s *Service) fetchInternal(ctx context.Context, text string, sources []string) ([]Result, error) {
	logger := ctx_pkg.L(ctx)

	// Валидация источников
	var targets []Provider
	for _, slug := range sources {
		if p, ok := s.providers[slug]; ok {
			targets = append(targets, p)
		} else {
			logger.Warn("unknown suggestion provider requested", slog.String("slug", slug))
		}
	}

	if len(targets) == 0 {
		return []Result{}, nil
	}

	// Канал для сбора результатов
	resultsCh := make(chan Result, len(targets))

	// Используем errgroup с контекстом.
	// Но! Мы не хотим отменять все запросы, если упал один.
	// Поэтому errgroup используем просто для Wait(), а ошибки обрабатываем локально.
	g, gCtx := errgroup.WithContext(ctx)

	for _, provider := range targets {
		p := provider // capture loop var
		g.Go(func() error {
			// У каждого провайдера должен быть свой таймаут, чтобы не вешать общий запрос надолго.
			// (Обычно это делается внутри клиента, но safety net здесь не помешает).
			childCtx, cancel := context.WithTimeout(gCtx, 5*time.Second)
			defer cancel()

			start := time.Now()
			res, err := p.Fetch(childCtx, text)
			duration := time.Since(start)

			if err != nil {
				// Логируем ошибку, но НЕ возвращаем её в errgroup, чтобы не отменить остальные.
				logger.Error("suggestion provider failed",
					slog.String("provider", p.Slug()),
					slog.String("text", text),
					slog.Duration("duration", duration),
					slog.Any("error", err),
				)
				return nil
			}

			// Если результат пустой (но без ошибки), тоже можно залогировать
			if res == nil {
				logger.Debug("suggestion provider returned no data", slog.String("provider", p.Slug()))
				return nil
			}

			select {
			case resultsCh <- *res:
			case <-gCtx.Done():
				return gCtx.Err()
			}

			return nil
		})
	}

	// Ждем завершения всех горутин
	if err := g.Wait(); err != nil {
		// Ошибка может возникнуть только если контекст отменен сверху
		return nil, err
	}
	close(resultsCh)

	// Собираем результаты
	var results []Result
	for res := range resultsCh {
		results = append(results, res)
	}

	return results, nil
}
