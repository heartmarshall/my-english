package repository

import (
	"context"
	"strings"
	"unicode/utf8"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
)

type DictionaryFilter struct {
	Search       string
	PartOfSpeech *model.PartOfSpeech
	HasCard      *bool

	Limit  int
	Offset int

	SortBy  *model.WordSortField
	SortDir *model.SortDirection
}

type DictionaryRepository struct {
	*Base[model.DictionaryEntry]
}

func NewDictionaryRepository(q database.Querier) *DictionaryRepository {
	return &DictionaryRepository{
		Base: NewBase[model.DictionaryEntry](
			q,
			schema.DictionaryEntries.Name.String(),
			schema.DictionaryEntries.Columns(),
		),
	}
}

func (r *DictionaryRepository) FindByNormalizedText(ctx context.Context, text string) (*model.DictionaryEntry, error) {
	return r.FindOneBy(ctx, schema.DictionaryEntries.TextNormalized.String(), text)
}

// applyFilters добавляет условия WHERE к билдеру
func (r *DictionaryRepository) applyFilters(b squirrel.SelectBuilder, f DictionaryFilter) squirrel.SelectBuilder {
	textCol := schema.DictionaryEntries.Text.String()
	idCol := schema.DictionaryEntries.ID.String()

	// 1. Фильтр по PartOfSpeech (через подзапрос EXISTS)
	if f.PartOfSpeech != nil {
		subQuery := squirrel.Select("1").
			From(schema.Senses.Name.String()).
			Where(squirrel.Eq{
				schema.Senses.EntryID.String():      squirrel.Expr(idCol),
				schema.Senses.PartOfSpeech.String(): *f.PartOfSpeech,
			})

		sql, args, _ := subQuery.ToSql()
		b = b.Where("EXISTS ("+sql+")", args...)
	}

	// 2. Фильтр по HasCard (через подзапрос EXISTS)
	if f.HasCard != nil {
		subQuery := squirrel.Select("1").
			From(schema.Cards.Name.String()).
			Where(squirrel.Eq{
				schema.Cards.EntryID.String(): squirrel.Expr(idCol),
			})

		sql, args, _ := subQuery.ToSql()
		if *f.HasCard {
			b = b.Where("EXISTS ("+sql+")", args...)
		} else {
			b = b.Where("NOT EXISTS ("+sql+")", args...)
		}
	}

	// 3. Поиск (Prefix или Trigrams)
	cleanQuery := strings.TrimSpace(f.Search)
	if cleanQuery != "" {
		queryLen := utf8.RuneCountInString(cleanQuery)
		if queryLen < 3 {
			// Prefix search for short words
			b = b.Where(squirrel.ILike{textCol: cleanQuery + "%"})
		} else {
			// Fuzzy search via pg_trgm for longer words
			b = b.Where(squirrel.Expr("? % ?", schema.DictionaryEntries.Text, cleanQuery))
		}
	}

	return b
}

// Find выполняет поиск слов
func (r *DictionaryRepository) Find(ctx context.Context, f DictionaryFilter) ([]model.DictionaryEntry, error) {
	b := r.SelectBuilder()
	b = r.applyFilters(b, f)

	// --- SORTING ---
	textCol := schema.DictionaryEntries.Text.String()
	cleanQuery := strings.TrimSpace(f.Search)
	appliedSort := false

	// A. Явная сортировка
	if f.SortBy != nil {
		direction := "ASC"
		if f.SortDir != nil && *f.SortDir == model.SortDirDesc {
			direction = "DESC"
		}

		switch *f.SortBy {
		case model.SortFieldText:
			b = b.OrderBy(textCol + " " + direction)
		case model.SortFieldCreatedAt:
			b = b.OrderBy(schema.DictionaryEntries.CreatedAt.String() + " " + direction)
		case model.SortFieldUpdatedAt:
			b = b.OrderBy(schema.DictionaryEntries.UpdatedAt.String() + " " + direction)
		}
		appliedSort = true
	}

	// B. Сортировка по релевантности (при поиске)
	if !appliedSort && cleanQuery != "" {
		queryLen := utf8.RuneCountInString(cleanQuery)
		if queryLen < 3 {
			// Короткие: сначала короткие совпадения, потом алфавит
			b = b.OrderBy("LENGTH(" + textCol + ") ASC")
			b = b.OrderBy(textCol + " ASC")
		} else {
			// Длинные: расстояние триграмм
			b = b.OrderByClause("? <-> ? ASC", schema.DictionaryEntries.Text, cleanQuery)
		}
		appliedSort = true
	}

	// C. Дефолтная сортировка
	if !appliedSort {
		b = b.OrderBy(schema.DictionaryEntries.CreatedAt.String() + " DESC")
	}

	// --- PAGINATION ---
	if f.Limit > 0 {
		b = b.Limit(uint64(f.Limit))
	}
	if f.Offset > 0 {
		b = b.Offset(uint64(f.Offset))
	}

	return r.List(ctx, b)
}

// CountTotal возвращает общее количество записей по фильтру (без учета лимитов)
func (r *DictionaryRepository) CountTotal(ctx context.Context, f DictionaryFilter) (int64, error) {
	b := squirrel.Select("COUNT(*)").From(schema.DictionaryEntries.Name.String())
	b = r.applyFilters(b, f)

	// Используем Querier напрямую, так как Base.Count не поддерживает наш специфичный applyFilters
	sql, args, err := b.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return 0, database.WrapDBError(err)
	}

	var count int64
	if err := r.Q().QueryRow(ctx, sql, args...).Scan(&count); err != nil {
		return 0, database.WrapDBError(err)
	}

	return count, nil
}

// Create создает новое слово
// Внимание: если слово с таким text_normalized уже существует, вернется ошибка.
// Для безопасного создания используйте CreateOrGet.
func (r *DictionaryRepository) Create(ctx context.Context, entry *model.DictionaryEntry) (*model.DictionaryEntry, error) {
	insert := r.InsertBuilder().
		Columns(schema.DictionaryEntries.InsertColumns()...).
		Values(entry.Text, entry.TextNormalized)

	return r.InsertReturning(ctx, insert)
}

// CreateOrGet создает новое слово или возвращает существующее по text_normalized
func (r *DictionaryRepository) CreateOrGet(ctx context.Context, entry *model.DictionaryEntry) (*model.DictionaryEntry, error) {
	// Сначала проверяем существование
	existing, err := r.FindByNormalizedText(ctx, entry.TextNormalized)
	if err == nil {
		return existing, nil
	}
	if err != database.ErrNotFound {
		return nil, err
	}

	// Слова нет, создаем новое
	return r.Create(ctx, entry)
}

// Update обновляет слово
func (r *DictionaryRepository) Update(ctx context.Context, id uuid.UUID, entry *model.DictionaryEntry) (*model.DictionaryEntry, error) {
	update := r.UpdateBuilder().
		Set("text", entry.Text).
		Set("text_normalized", entry.TextNormalized).
		Where(squirrel.Eq{schema.DictionaryEntries.ID.String(): id}).
		Suffix("RETURNING *")

	sql, args, err := update.ToSql()
	if err != nil {
		return nil, database.WrapDBError(err)
	}

	var result model.DictionaryEntry
	if err := pgxscan.Get(ctx, r.Q(), &result, sql, args...); err != nil {
		if pgxscan.NotFound(err) {
			return nil, database.ErrNotFound
		}
		return nil, database.WrapDBError(err)
	}
	return &result, nil
}

// Delete удаляет слово (CASCADE удалит связанные данные)
func (r *DictionaryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.Base.Delete(ctx, schema.DictionaryEntries.ID.String(), id)
}

// ExistsByNormalizedText проверяет существование слова по нормализованному тексту
func (r *DictionaryRepository) ExistsByNormalizedText(ctx context.Context, text string) (bool, error) {
	query := r.SelectBuilder().
		Columns("1").
		Where(squirrel.Eq{schema.DictionaryEntries.TextNormalized.String(): text}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return false, database.WrapDBError(err)
	}

	var exists int
	err = r.Q().QueryRow(ctx, sql, args...).Scan(&exists)
	if err != nil {
		if pgxscan.NotFound(err) {
			return false, nil
		}
		return false, database.WrapDBError(err)
	}
	return true, nil
}
