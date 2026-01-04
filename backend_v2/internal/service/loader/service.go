package loader

import (
	"context"

	"github.com/google/uuid"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/repository"
	factory "github.com/heartmarshall/my-english/internal/database/repository/factory"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

type Service struct {
	repos *factory.Factory
	db    database.Querier // Пул соединений
}

func New(repos *factory.Factory, db database.Querier) *Service {
	return &Service{
		repos: repos,
		db:    db,
	}
}

// --- Linguistic ---

func (s *Service) GetSensesByLexemeIDs(ctx context.Context, lexemeIDs []uuid.UUID) ([]model.Sense, error) {
	return s.repos.Sense(s.db).ListByLexemeIDs(ctx, lexemeIDs)
}

func (s *Service) GetPronunciationsByLexemeIDs(ctx context.Context, lexemeIDs []uuid.UUID) ([]model.Pronunciation, error) {
	return s.repos.Pronunciation(s.db).ListByLexemeIDs(ctx, lexemeIDs)
}

func (s *Service) GetTranslationsBySenseIDs(ctx context.Context, senseIDs []uuid.UUID) ([]model.SenseTranslation, error) {
	return s.repos.SenseTranslation(s.db).ListBySenseIDs(ctx, senseIDs)
}

func (s *Service) GetExamplesBySenseIDs(ctx context.Context, senseIDs []uuid.UUID) ([]model.Example, error) {
	return s.repos.Example(s.db).ListBySenseIDs(ctx, senseIDs)
}

// --- User / Study ---

func (s *Service) GetSRSByCardIDs(ctx context.Context, cardIDs []uuid.UUID) ([]model.SRSState, error) {
	return s.repos.SRS(s.db).ListByCardIDs(ctx, cardIDs)
}

// GetTagsByCardIDs — сложный случай (Many-to-Many).
// Возвращает мапу: CardID -> []Tag
func (s *Service) GetTagsByCardIDs(ctx context.Context, cardIDs []uuid.UUID) (map[uuid.UUID][]model.Tag, error) {
	// 1. Получаем связи (CardID, TagID)
	links, err := s.repos.CardTag(s.db).ListByCardIDs(ctx, cardIDs)
	if err != nil {
		return nil, err
	}

	if len(links) == 0 {
		return map[uuid.UUID][]model.Tag{}, nil
	}

	// 2. Собираем уникальные TagID
	tagIDsMap := make(map[int]bool)
	for _, link := range links {
		tagIDsMap[link.TagID] = true
	}

	uniqueTagIDs := make([]int, 0, len(tagIDsMap))
	for id := range tagIDsMap {
		uniqueTagIDs = append(uniqueTagIDs, id)
	}

	// 3. Загружаем сами теги
	tags, err := s.repos.Tag(s.db).GetByIDs(ctx, uniqueTagIDs)
	if err != nil {
		return nil, err
	}

	// 4. Создаем мапу TagID -> Tag для быстрого поиска
	tagsByID := make(map[int]model.Tag)
	for _, t := range tags {
		tagsByID[t.ID] = t
	}

	// 5. Собираем итоговую мапу CardID -> []Tag
	result := make(map[uuid.UUID][]model.Tag)
	for _, link := range links {
		if tag, ok := tagsByID[link.TagID]; ok {
			result[link.CardID] = append(result[link.CardID], tag)
		}
	}

	return result, nil
}

// GetSensesByIDs возвращает список смыслов по списку ID.
func (s *Service) GetSensesByIDs(ctx context.Context, senseIDs []uuid.UUID) ([]model.Sense, error) {
	return s.repos.Sense(s.db).GetByIDs(ctx, senseIDs)
}

// GetLexemesByIDs возвращает список лексем по списку ID.
func (s *Service) GetLexemesByIDs(ctx context.Context, lexemeIDs []uuid.UUID) ([]model.Lexeme, error) {
	if len(lexemeIDs) == 0 {
		return []model.Lexeme{}, nil
	}
	return s.repos.Lexeme(s.db).FindBy(ctx, schema.Lexemes.ID.String(), lexemeIDs)
}

// GetDataSourcesByIDs возвращает список источников данных по списку ID.
func (s *Service) GetDataSourcesByIDs(ctx context.Context, sourceIDs []int) ([]model.DataSource, error) {
	if len(sourceIDs) == 0 {
		return []model.DataSource{}, nil
	}
	// Преобразуем []int в []any для FindBy
	ids := make([]any, len(sourceIDs))
	for i, id := range sourceIDs {
		ids[i] = id
	}
	return s.repos.DataSource(s.db).FindBy(ctx, schema.DataSources.ID.String(), ids)
}

// GetRelationsBySenseIDs возвращает список связей для списка смыслов.
func (s *Service) GetRelationsBySenseIDs(ctx context.Context, senseIDs []uuid.UUID) ([]model.SenseRelation, error) {
	return s.repos.SenseRelation(s.db).ListBySenseIDs(ctx, senseIDs)
}

// GetInflectionsByLemmaID возвращает формы для леммы.
func (s *Service) GetInflectionsByLemmaID(ctx context.Context, lemmaID uuid.UUID) ([]model.Inflection, error) {
	return s.repos.Inflection(s.db).GetFormsByLemmaID(ctx, lemmaID)
}

// GetLemmaByInflectedID возвращает лемму для формы.
func (s *Service) GetLemmaByInflectedID(ctx context.Context, inflectedID uuid.UUID) (*model.Inflection, error) {
	inflections, err := s.repos.Inflection(s.db).GetLemmaByInflectedID(ctx, inflectedID)
	if err != nil {
		return nil, err
	}
	if len(inflections) == 0 {
		return nil, nil
	}
	return &inflections[0], nil
}

// GetReviewCountByCardIDs возвращает количество повторений для карточек.
func (s *Service) GetReviewCountByCardIDs(ctx context.Context, cardIDs []uuid.UUID) (map[uuid.UUID]int, error) {
	result := make(map[uuid.UUID]int)

	// Для каждой карточки считаем количество review logs
	for _, cardID := range cardIDs {
		// Используем Count с условием через schema
		count, err := s.repos.Review(s.db).Count(ctx, repository.WithWhere(schema.ReviewLogs.CardID.Eq(cardID)))
		if err != nil {
			return nil, err
		}
		result[cardID] = int(count)
	}

	return result, nil
}
