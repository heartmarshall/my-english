package meaningtag

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/model"
)

const (
	tableName = "meanings_tags"
)

// Repo — реализация репозитория для работы со связями meaning-tag.
type Repo struct {
	q database.Querier
}

// New создаёт новый репозиторий.
func New(q database.Querier) *Repo {
	return &Repo{q: q}
}

// AttachTag привязывает tag к meaning.
func (r *Repo) AttachTag(ctx context.Context, meaningID, tagID int64) error {
	query, args, err := database.Builder.
		Insert(tableName).
		Columns("meaning_id", "tag_id").
		Values(meaningID, tagID).
		Suffix("ON CONFLICT DO NOTHING").
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.q.Exec(ctx, query, args...)
	return err
}

// AttachTags привязывает несколько tags к meaning.
func (r *Repo) AttachTags(ctx context.Context, meaningID int64, tagIDs []int64) error {
	if len(tagIDs) == 0 {
		return nil
	}

	qb := database.Builder.
		Insert(tableName).
		Columns("meaning_id", "tag_id")

	for _, tagID := range tagIDs {
		qb = qb.Values(meaningID, tagID)
	}

	qb = qb.Suffix("ON CONFLICT DO NOTHING")

	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	_, err = r.q.Exec(ctx, query, args...)
	return err
}

// DetachTag отвязывает tag от meaning.
func (r *Repo) DetachTag(ctx context.Context, meaningID, tagID int64) error {
	query, args, err := database.Builder.
		Delete(tableName).
		Where(squirrel.Eq{"meaning_id": meaningID, "tag_id": tagID}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.q.Exec(ctx, query, args...)
	return err
}

// DetachAllFromMeaning отвязывает все tags от meaning.
func (r *Repo) DetachAllFromMeaning(ctx context.Context, meaningID int64) error {
	query, args, err := database.Builder.
		Delete(tableName).
		Where(squirrel.Eq{"meaning_id": meaningID}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.q.Exec(ctx, query, args...)
	return err
}

// GetTagIDsByMeaningID возвращает ID тегов для meaning.
func (r *Repo) GetTagIDsByMeaningID(ctx context.Context, meaningID int64) ([]int64, error) {
	query, args, err := database.Builder.
		Select("tag_id").
		From(tableName).
		Where(squirrel.Eq{"meaning_id": meaningID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.q.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make([]int64, 0)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, rows.Err()
}

// GetMeaningIDsByTagID возвращает ID meanings для tag.
func (r *Repo) GetMeaningIDsByTagID(ctx context.Context, tagID int64) ([]int64, error) {
	query, args, err := database.Builder.
		Select("meaning_id").
		From(tableName).
		Where(squirrel.Eq{"tag_id": tagID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.q.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make([]int64, 0)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, rows.Err()
}

// GetByMeaningIDs возвращает все связи для нескольких meanings (для batch loading).
func (r *Repo) GetByMeaningIDs(ctx context.Context, meaningIDs []int64) ([]*model.MeaningTag, error) {
	if len(meaningIDs) == 0 {
		return make([]*model.MeaningTag, 0), nil
	}

	query, args, err := database.Builder.
		Select("meaning_id", "tag_id").
		From(tableName).
		Where(squirrel.Eq{"meaning_id": meaningIDs}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.q.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*model.MeaningTag, 0)
	for rows.Next() {
		var mt model.MeaningTag
		if err := rows.Scan(&mt.MeaningID, &mt.TagID); err != nil {
			return nil, err
		}
		result = append(result, &mt)
	}

	return result, rows.Err()
}

// SyncTags синхронизирует теги meaning: удаляет старые и добавляет новые.
func (r *Repo) SyncTags(ctx context.Context, meaningID int64, tagIDs []int64) error {
	// Удаляем все старые связи
	if err := r.DetachAllFromMeaning(ctx, meaningID); err != nil {
		return err
	}

	// Добавляем новые
	return r.AttachTags(ctx, meaningID, tagIDs)
}
