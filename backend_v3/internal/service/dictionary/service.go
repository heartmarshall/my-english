package dictionary

import (
	"context"
	"errors"
	"strings"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/repository"
	"github.com/heartmarshall/my-english/internal/model"
	"github.com/heartmarshall/my-english/internal/service/types"
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

// CreateWord создает слово и все связанные сущности атомарно.
func (s *Service) CreateWord(ctx context.Context, input CreateWordInput) (*model.DictionaryEntry, error) {
	// 1. Нормализация и валидация
	textRaw := strings.TrimSpace(input.Text)
	textNorm := normalizeText(textRaw)

	if textRaw == "" {
		return nil, errors.New("text cannot be empty")
	}

	var createdEntry *model.DictionaryEntry

	// 2. Транзакция
	err := s.tx.RunInTx(ctx, func(ctx context.Context, _ database.Querier) error {
		// A. Проверяем дубликат
		exists, err := s.repos.Dictionary.ExistsByNormalizedText(ctx, textNorm)
		if err != nil {
			return err
		}
		if exists {
			return types.ErrAlreadyExists
		}

		// B. Создаем основную запись (Entry)
		entry := &model.DictionaryEntry{
			Text:           textRaw,
			TextNormalized: textNorm,
		}
		createdEntry, err = s.repos.Dictionary.Create(ctx, entry)
		if err != nil {
			return err
		}

		// C. Создаем смыслы (Senses) и вложенные в них сущности
		// Примечание: Insert Senses делаем в цикле, так как нам нужны ID созданных смыслов
		// для вставки переводов и примеров.
		for _, senseIn := range input.Senses {
			sense := &model.Sense{
				EntryID:      createdEntry.ID,
				Definition:   senseIn.Definition,
				PartOfSpeech: senseIn.PartOfSpeech,
				SourceSlug:   senseIn.SourceSlug,
			}

			createdSense, err := s.repos.Senses.Create(ctx, sense)
			if err != nil {
				return err
			}

			// C.1 Переводы (Translations) - Batch Insert
			if len(senseIn.Translations) > 0 {
				translations := make([]model.Translation, len(senseIn.Translations))
				for i, tr := range senseIn.Translations {
					translations[i] = model.Translation{
						SenseID:    createdSense.ID,
						Text:       tr.Text,
						SourceSlug: tr.SourceSlug,
					}
				}
				if _, err := s.repos.Translations.BatchCreate(ctx, translations); err != nil {
					return err
				}
			}

			// C.2 Примеры (Examples) - Batch Insert
			if len(senseIn.Examples) > 0 {
				examples := make([]model.Example, len(senseIn.Examples))
				for i, ex := range senseIn.Examples {
					examples[i] = model.Example{
						SenseID:     createdSense.ID,
						Sentence:    ex.Sentence,
						Translation: ex.Translation,
						SourceSlug:  ex.SourceSlug,
					}
				}
				if _, err := s.repos.Examples.BatchCreate(ctx, examples); err != nil {
					return err
				}
			}
		}

		// D. Изображения (Images) - Batch Insert
		if len(input.Images) > 0 {
			images := make([]model.Image, len(input.Images))
			for i, img := range input.Images {
				images[i] = model.Image{
					EntryID:    createdEntry.ID,
					URL:        img.URL,
					Caption:    img.Caption,
					SourceSlug: img.SourceSlug,
				}
			}
			if _, err := s.repos.Images.BatchCreate(ctx, images); err != nil {
				return err
			}
		}

		// E. Произношения (Pronunciations) - Batch Insert
		if len(input.Pronunciations) > 0 {
			prons := make([]model.Pronunciation, len(input.Pronunciations))
			for i, p := range input.Pronunciations {
				prons[i] = model.Pronunciation{
					EntryID:       createdEntry.ID,
					AudioURL:      p.AudioURL,
					Transcription: p.Transcription,
					Region:        p.Region,
					SourceSlug:    p.SourceSlug,
				}
			}
			if _, err := s.repos.Pronunciations.BatchCreate(ctx, prons); err != nil {
				return err
			}
		}

		// F. Карточка для изучения (Card) - Опционально
		if input.CreateCard {
			card := &model.Card{
				EntryID:      createdEntry.ID,
				Status:       model.StatusNew,
				IntervalDays: 0,
				EaseFactor:   2.5, // Дефолтное значение для SM-2
			}
			if _, err := s.repos.Cards.Create(ctx, card); err != nil {
				return err
			}
		}

		// G. Audit Log (Создание записи)
		// Используем JSON для сохранения деталей (снапшот не делаем, просто факт создания)
		audit := &model.AuditRecord{
			EntityType: model.EntityEntry,
			EntityID:   &createdEntry.ID,
			Action:     model.ActionCreate,
			Changes:    model.JSON{"text": createdEntry.Text},
		}
		if _, err := s.repos.Audit.Create(ctx, audit); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		// Маппинг ошибок БД на доменные ошибки, если нужно уточнить
		if database.IsDuplicateError(err) {
			return nil, types.ErrAlreadyExists
		}
		return nil, err
	}

	return createdEntry, nil
}

// Helper для нормализации
func normalizeText(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
