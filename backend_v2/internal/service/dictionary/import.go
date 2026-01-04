package dictionary

import (
	"context"
	"fmt"
	"strings"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/model"
)

// SaveImportedWord сохраняет импортированное слово и все связанные данные (смыслы, переводы, примеры)
// в базу данных в рамках одной транзакции.
//
// sourceSlug — строковый идентификатор источника (например, "freedict"), который должен существовать в таблице data_sources.
func (s *Service) SaveImportedWord(ctx context.Context, word *ImportedWord, sourceSlug string) (*model.Lexeme, error) {
	var lexeme *model.Lexeme

	// Запускаем транзакцию. TxManager создаст tx и передаст её в функцию.
	err := s.txManager.RunInTx(ctx, func(ctx context.Context, tx database.Querier) error {
		// 1. Получаем ID источника данных по слагу
		// Это важно для FK в таблицах senses, pronunciations и т.д.
		dataSource, err := s.repos.DataSource(tx).GetBySlug(ctx, sourceSlug)
		if err != nil {
			return fmt.Errorf("failed to get data source '%s': %w", sourceSlug, err)
		}

		// 2. Создаем или находим Лексему (Слово)
		// Нормализуем текст (нижний регистр, trim) для поиска/создания
		normalizedText := strings.TrimSpace(strings.ToLower(word.Text))

		inputLexeme := &model.Lexeme{
			TextNormalized: normalizedText,
			TextDisplay:    strings.TrimSpace(word.Text),
		}

		// Используем репозиторий с внедренной транзакцией tx
		lexemeRepo := s.repos.Lexeme(tx)
		lexeme, err = lexemeRepo.CreateWithConflictIgnore(ctx, inputLexeme)
		if err != nil {
			return fmt.Errorf("failed to upsert lexeme: %w", err)
		}

		// 3. Сохраняем варианты произношения
		if len(word.Pronunciations) > 0 {
			pronRepo := s.repos.Pronunciation(tx)
			for _, p := range word.Pronunciations {
				// Создаем запись произношения
				_, err := pronRepo.Create(ctx, &model.Pronunciation{
					LexemeID:      lexeme.ID,
					AudioURL:      p.AudioURL,
					Transcription: &p.Transcription,
					Region:        string(p.Region),
					SourceID:      &dataSource.ID,
				})
				if err != nil {
					return fmt.Errorf("failed to create pronunciation: %w", err)
				}
			}
		}

		// 4. Сохраняем Смыслы (Senses) и вложенные в них данные
		senseRepo := s.repos.Sense(tx)
		transRepo := s.repos.SenseTranslation(tx)
		exRepo := s.repos.Example(tx)

		for _, senseData := range word.Senses {
			// 4.1 Создаем Sense
			sense, err := senseRepo.Create(ctx, &model.Sense{
				LexemeID:     lexeme.ID,
				PartOfSpeech: string(senseData.PartOfSpeech),
				Definition:   senseData.Definition,
				SourceID:     dataSource.ID,
				// CefrLevel и ExternalRefID можно добавить в DTO позже, если API их отдает
			})
			if err != nil {
				return fmt.Errorf("failed to create sense: %w", err)
			}

			// 4.2 Добавляем переводы к этому смыслу
			for _, tr := range senseData.Translations {
				if tr == "" {
					continue
				}
				_, err := transRepo.Create(ctx, &model.SenseTranslation{
					SenseID:     sense.ID,
					Translation: tr,
					SourceID:    &dataSource.ID,
				})
				if err != nil {
					return fmt.Errorf("failed to create translation: %w", err)
				}
			}

			// 4.3 Добавляем примеры к этому смыслу
			for _, ex := range senseData.Examples {
				// Пропускаем пустые примеры
				if ex.SentenceEn == "" {
					continue
				}

				exampleModel := &model.Example{
					SenseID:    &sense.ID,
					SentenceEn: ex.SentenceEn,
				}

				if ex.SentenceRu != "" {
					exampleModel.SentenceRu = &ex.SentenceRu
				}

				// SourceName пока оставляем пустым или можно брать DisplayName источника

				_, err := exRepo.Create(ctx, exampleModel)
				if err != nil {
					return fmt.Errorf("failed to create example: %w", err)
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return lexeme, nil
}
