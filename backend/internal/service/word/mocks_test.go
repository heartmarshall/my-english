package word_test

import (
	"context"

	"github.com/heartmarshall/my-english/internal/model"
)

// --- Mock implementations ---

type mockWordRepository struct {
	CreateFunc    func(ctx context.Context, word *model.Word) error
	GetByIDFunc   func(ctx context.Context, id int64) (*model.Word, error)
	GetByTextFunc func(ctx context.Context, text string) (*model.Word, error)
	ListFunc      func(ctx context.Context, filter *model.WordFilter, limit, offset int) ([]*model.Word, error)
	CountFunc     func(ctx context.Context, filter *model.WordFilter) (int, error)
	UpdateFunc    func(ctx context.Context, word *model.Word) error
	DeleteFunc    func(ctx context.Context, id int64) error
}

func (m *mockWordRepository) Create(ctx context.Context, word *model.Word) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, word)
	}
	return nil
}

func (m *mockWordRepository) GetByID(ctx context.Context, id int64) (*model.Word, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockWordRepository) GetByText(ctx context.Context, text string) (*model.Word, error) {
	if m.GetByTextFunc != nil {
		return m.GetByTextFunc(ctx, text)
	}
	return nil, nil
}

func (m *mockWordRepository) List(ctx context.Context, filter *model.WordFilter, limit, offset int) ([]*model.Word, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, filter, limit, offset)
	}
	return nil, nil
}

func (m *mockWordRepository) Count(ctx context.Context, filter *model.WordFilter) (int, error) {
	if m.CountFunc != nil {
		return m.CountFunc(ctx, filter)
	}
	return 0, nil
}

func (m *mockWordRepository) Update(ctx context.Context, word *model.Word) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, word)
	}
	return nil
}

func (m *mockWordRepository) Delete(ctx context.Context, id int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

type mockMeaningRepository struct {
	CreateFunc         func(ctx context.Context, meaning *model.Meaning) error
	GetByIDFunc        func(ctx context.Context, id int64) (*model.Meaning, error)
	GetByWordIDFunc    func(ctx context.Context, wordID int64) ([]*model.Meaning, error)
	UpdateFunc         func(ctx context.Context, meaning *model.Meaning) error
	DeleteFunc         func(ctx context.Context, id int64) error
	DeleteByWordIDFunc func(ctx context.Context, wordID int64) (int64, error)
}

func (m *mockMeaningRepository) Create(ctx context.Context, meaning *model.Meaning) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, meaning)
	}
	meaning.ID = 1
	return nil
}

func (m *mockMeaningRepository) GetByID(ctx context.Context, id int64) (*model.Meaning, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockMeaningRepository) GetByWordID(ctx context.Context, wordID int64) ([]*model.Meaning, error) {
	if m.GetByWordIDFunc != nil {
		return m.GetByWordIDFunc(ctx, wordID)
	}
	return []*model.Meaning{}, nil
}

func (m *mockMeaningRepository) Update(ctx context.Context, meaning *model.Meaning) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, meaning)
	}
	return nil
}

func (m *mockMeaningRepository) Delete(ctx context.Context, id int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *mockMeaningRepository) DeleteByWordID(ctx context.Context, wordID int64) (int64, error) {
	if m.DeleteByWordIDFunc != nil {
		return m.DeleteByWordIDFunc(ctx, wordID)
	}
	return 0, nil
}

type mockExampleRepository struct {
	CreateFunc            func(ctx context.Context, example *model.Example) error
	CreateBatchFunc       func(ctx context.Context, examples []*model.Example) error
	GetByMeaningIDFunc    func(ctx context.Context, meaningID int64) ([]*model.Example, error)
	GetByMeaningIDsFunc   func(ctx context.Context, meaningIDs []int64) ([]*model.Example, error)
	DeleteByMeaningIDFunc func(ctx context.Context, meaningID int64) (int64, error)
}

func (m *mockExampleRepository) Create(ctx context.Context, example *model.Example) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, example)
	}
	return nil
}

func (m *mockExampleRepository) CreateBatch(ctx context.Context, examples []*model.Example) error {
	if m.CreateBatchFunc != nil {
		return m.CreateBatchFunc(ctx, examples)
	}
	for i := range examples {
		examples[i].ID = int64(i + 1)
	}
	return nil
}

func (m *mockExampleRepository) GetByMeaningID(ctx context.Context, meaningID int64) ([]*model.Example, error) {
	if m.GetByMeaningIDFunc != nil {
		return m.GetByMeaningIDFunc(ctx, meaningID)
	}
	return []*model.Example{}, nil
}

func (m *mockExampleRepository) GetByMeaningIDs(ctx context.Context, meaningIDs []int64) ([]*model.Example, error) {
	if m.GetByMeaningIDsFunc != nil {
		return m.GetByMeaningIDsFunc(ctx, meaningIDs)
	}
	return []*model.Example{}, nil
}

func (m *mockExampleRepository) DeleteByMeaningID(ctx context.Context, meaningID int64) (int64, error) {
	if m.DeleteByMeaningIDFunc != nil {
		return m.DeleteByMeaningIDFunc(ctx, meaningID)
	}
	return 0, nil
}

type mockTagRepository struct {
	GetByNameFunc   func(ctx context.Context, name string) (*model.Tag, error)
	GetByNamesFunc  func(ctx context.Context, names []string) ([]*model.Tag, error)
	GetByIDsFunc    func(ctx context.Context, ids []int64) ([]*model.Tag, error)
	GetOrCreateFunc func(ctx context.Context, name string) (*model.Tag, error)
}

func (m *mockTagRepository) GetByName(ctx context.Context, name string) (*model.Tag, error) {
	if m.GetByNameFunc != nil {
		return m.GetByNameFunc(ctx, name)
	}
	return nil, nil
}

func (m *mockTagRepository) GetByNames(ctx context.Context, names []string) ([]*model.Tag, error) {
	if m.GetByNamesFunc != nil {
		return m.GetByNamesFunc(ctx, names)
	}
	return []*model.Tag{}, nil
}

func (m *mockTagRepository) GetByIDs(ctx context.Context, ids []int64) ([]*model.Tag, error) {
	if m.GetByIDsFunc != nil {
		return m.GetByIDsFunc(ctx, ids)
	}
	return []*model.Tag{}, nil
}

func (m *mockTagRepository) GetOrCreate(ctx context.Context, name string) (*model.Tag, error) {
	if m.GetOrCreateFunc != nil {
		return m.GetOrCreateFunc(ctx, name)
	}
	return &model.Tag{ID: 1, Name: name}, nil
}

type mockMeaningTagRepository struct {
	AttachTagsFunc           func(ctx context.Context, meaningID int64, tagIDs []int64) error
	GetTagIDsByMeaningIDFunc func(ctx context.Context, meaningID int64) ([]int64, error)
	GetByMeaningIDsFunc      func(ctx context.Context, meaningIDs []int64) ([]*model.MeaningTag, error)
	SyncTagsFunc             func(ctx context.Context, meaningID int64, tagIDs []int64) error
	DetachAllFromMeaningFunc func(ctx context.Context, meaningID int64) error
}

func (m *mockMeaningTagRepository) AttachTags(ctx context.Context, meaningID int64, tagIDs []int64) error {
	if m.AttachTagsFunc != nil {
		return m.AttachTagsFunc(ctx, meaningID, tagIDs)
	}
	return nil
}

func (m *mockMeaningTagRepository) GetTagIDsByMeaningID(ctx context.Context, meaningID int64) ([]int64, error) {
	if m.GetTagIDsByMeaningIDFunc != nil {
		return m.GetTagIDsByMeaningIDFunc(ctx, meaningID)
	}
	return []int64{}, nil
}

func (m *mockMeaningTagRepository) GetByMeaningIDs(ctx context.Context, meaningIDs []int64) ([]*model.MeaningTag, error) {
	if m.GetByMeaningIDsFunc != nil {
		return m.GetByMeaningIDsFunc(ctx, meaningIDs)
	}
	return []*model.MeaningTag{}, nil
}

func (m *mockMeaningTagRepository) SyncTags(ctx context.Context, meaningID int64, tagIDs []int64) error {
	if m.SyncTagsFunc != nil {
		return m.SyncTagsFunc(ctx, meaningID, tagIDs)
	}
	return nil
}

func (m *mockMeaningTagRepository) DetachAllFromMeaning(ctx context.Context, meaningID int64) error {
	if m.DetachAllFromMeaningFunc != nil {
		return m.DetachAllFromMeaningFunc(ctx, meaningID)
	}
	return nil
}
