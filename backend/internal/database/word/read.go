package word

import (
	"context"
	"log/slog"
	"strconv"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/database/schema"
	"github.com/heartmarshall/my-english/internal/model"
	ctxlog "github.com/heartmarshall/my-english/pkg/context"
)

func (r *Repo) GetByID(ctx context.Context, id int64) (model.Word, error) {
	builder := database.Builder.
		Select(schema.Words.All()...).
		From(schema.Words.Name.String()).
		Where(schema.Words.ID.Eq(id))

	return database.NewQuery[model.Word](r.q, builder).One(ctx)
}

func (r *Repo) GetByText(ctx context.Context, text string) (model.Word, error) {
	builder := database.Builder.
		Select(schema.Words.All()...).
		From(schema.Words.Name.String()).
		Where(schema.Words.Text.Eq(text))

	return database.NewQuery[model.Word](r.q, builder).One(ctx)
}

func (r *Repo) List(ctx context.Context, filter *model.WordFilter, limit, offset int) ([]model.Word, error) {
	limit, offset = database.NormalizePagination(limit, offset)

	// Если есть поиск, используем триграммный поиск с прямой SQL
	if filter != nil && filter.Search != nil && *filter.Search != "" {
		return r.listWithTrigramSearch(ctx, filter, limit, offset)
	}

	// Обычный поиск без триграмм
	selectCols := make([]string, 0, len(schema.Words.All()))
	for _, col := range schema.Words.All() {
		selectCols = append(selectCols, string(col))
	}

	qb := database.Builder.
		Select(selectCols...).
		Distinct().
		From(schema.Words.Name.String())

	qb = applyFilter(qb, filter)

	qb = qb.
		OrderBy(schema.Words.CreatedAt.Desc()).
		Limit(uint64(limit)).
		Offset(uint64(offset))

	return database.NewQuery[model.Word](r.q, qb).List(ctx)
}

// listWithTrigramSearch выполняет поиск с использованием триграмм
func (r *Repo) listWithTrigramSearch(ctx context.Context, filter *model.WordFilter, limit, offset int) ([]model.Word, error) {
	searchQuery := *filter.Search
	similarityThreshold := 0.2

	// Формируем список колонок С префиксом таблицы для подзапроса
	selectColsForSubquery := make([]string, 0, len(schema.Words.All())+1)
	for _, col := range schema.Words.All() {
		selectColsForSubquery = append(selectColsForSubquery, string(col))
	}
	selectColsForSubquery = append(selectColsForSubquery, "word_similarity($1, "+string(schema.Words.Text)+") AS similarity")

	// Формируем список колонок БЕЗ префикса таблицы для финального SELECT
	// В финальном SELECT после подзапроса нужно использовать колонки без префикса таблицы
	selectColsFinal := make([]string, 0, len(schema.Words.All()))
	for _, col := range schema.Words.All() {
		// Убираем префикс таблицы (например, "words.id" -> "id")
		colStr := string(col)
		if idx := strings.Index(colStr, "."); idx != -1 {
			colStr = colStr[idx+1:]
		}
		selectColsFinal = append(selectColsFinal, colStr)
	}

	// Логируем для отладки
	ctxlog.L(ctx).Info("Preparing final SELECT columns",
		slog.Any("selectColsFinal", selectColsFinal),
		slog.Int("count", len(selectColsFinal)),
	)

	// Начинаем формировать SQL запрос (в подзапросе используем колонки С similarity и префиксом)
	var query strings.Builder
	query.WriteString("SELECT ")
	query.WriteString(strings.Join(selectColsForSubquery, ", "))
	query.WriteString(" FROM ")
	query.WriteString(schema.Words.Name.String())

	args := []interface{}{searchQuery}
	argIndex := 2 // $1 уже используется для searchQuery

	// Добавляем WHERE условия для триграммного поиска
	whereConditions := []string{
		"word_similarity($1, " + string(schema.Words.Text) + ") > $" + strconv.Itoa(argIndex) + " OR " +
			string(schema.Words.Text) + " % $1 OR " +
			string(schema.Words.Text) + " ILIKE $" + strconv.Itoa(argIndex+1),
	}
	args = append(args, similarityThreshold, "%"+searchQuery+"%")
	argIndex += 2

	// Добавляем JOIN и условия для статуса
	if filter.Status != nil {
		query.WriteString(" JOIN ")
		query.WriteString(schema.Meanings.Name.String())
		query.WriteString(" ON ")
		query.WriteString(string(schema.Meanings.WordID))
		query.WriteString(" = ")
		query.WriteString(string(schema.Words.ID))
		whereConditions = append(whereConditions, string(schema.Meanings.LearningStatus)+" = $"+strconv.Itoa(argIndex))
		args = append(args, *filter.Status)
		argIndex++
	}

	// Добавляем JOIN и условия для тегов
	if len(filter.Tags) > 0 {
		if filter.Status == nil {
			query.WriteString(" JOIN ")
			query.WriteString(schema.Meanings.Name.String())
			query.WriteString(" ON ")
			query.WriteString(string(schema.Meanings.WordID))
			query.WriteString(" = ")
			query.WriteString(string(schema.Words.ID))
		}
		query.WriteString(" JOIN ")
		query.WriteString(schema.MeaningTags.Name.String())
		query.WriteString(" ON ")
		query.WriteString(string(schema.MeaningTags.MeaningID))
		query.WriteString(" = ")
		query.WriteString(string(schema.Meanings.ID))
		query.WriteString(" JOIN ")
		query.WriteString(schema.Tags.Name.String())
		query.WriteString(" ON ")
		query.WriteString(string(schema.Tags.ID))
		query.WriteString(" = ")
		query.WriteString(string(schema.MeaningTags.TagID))

		// Формируем IN условие для тегов
		tagPlaceholders := make([]string, len(filter.Tags))
		for i := range filter.Tags {
			tagPlaceholders[i] = "$" + strconv.Itoa(argIndex)
			args = append(args, filter.Tags[i])
			argIndex++
		}
		whereConditions = append(whereConditions, string(schema.Tags.NameCol)+" IN ("+strings.Join(tagPlaceholders, ", ")+")")
	}

	// Объединяем WHERE условия
	if len(whereConditions) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(strings.Join(whereConditions, " AND "))
	}

	// Всегда оборачиваем в подзапрос, чтобы убрать similarity из финального SELECT
	// similarity используется только для сортировки в подзапросе
	// selectColsFinal уже не содержит similarity

	baseQuery := query.String()
	query.Reset()

	// Добавляем сортировку в подзапрос (здесь similarity еще доступен)
	baseQueryWithOrder := baseQuery + " ORDER BY similarity DESC, " + string(schema.Words.CreatedAt) + " DESC"

	// Формируем финальный запрос без similarity в SELECT
	needsDistinct := filter.Status != nil || len(filter.Tags) > 0
	if needsDistinct {
		// Для DISTINCT ON тоже нужно использовать колонку без префикса
		idColForDistinct := string(schema.Words.ID)
		if idx := strings.Index(idColForDistinct, "."); idx != -1 {
			idColForDistinct = idColForDistinct[idx+1:]
		}
		query.WriteString("SELECT DISTINCT ON (")
		query.WriteString(idColForDistinct)
		query.WriteString(") ")
		query.WriteString(strings.Join(selectColsFinal, ", "))
		query.WriteString(" FROM (")
		query.WriteString(baseQueryWithOrder)
		query.WriteString(") AS subquery")
		// Сортируем по ID, так как similarity уже использован в подзапросе
		// В ORDER BY тоже нужно использовать колонку без префикса таблицы
		idCol := string(schema.Words.ID)
		if idx := strings.Index(idCol, "."); idx != -1 {
			idCol = idCol[idx+1:]
		}
		query.WriteString(" ORDER BY ")
		query.WriteString(idCol)
		query.WriteString(" LIMIT $")
		query.WriteString(strconv.Itoa(argIndex))
		query.WriteString(" OFFSET $")
		query.WriteString(strconv.Itoa(argIndex + 1))
	} else {
		// Логируем перед формированием финального запроса
		ctxlog.L(ctx).Info("Building final SELECT without similarity",
			slog.String("selectColsFinal", strings.Join(selectColsFinal, ", ")),
		)
		query.WriteString("SELECT ")
		query.WriteString(strings.Join(selectColsFinal, ", "))
		query.WriteString(" FROM (")
		query.WriteString(baseQueryWithOrder)
		query.WriteString(") AS subquery")
		// Сортируем по ID, так как similarity уже использован в подзапросе
		// В ORDER BY тоже нужно использовать колонку без префикса таблицы
		idCol := string(schema.Words.ID)
		if idx := strings.Index(idCol, "."); idx != -1 {
			idCol = idCol[idx+1:]
		}
		query.WriteString(" ORDER BY ")
		query.WriteString(idCol)
		query.WriteString(" LIMIT $")
		query.WriteString(strconv.Itoa(argIndex))
		query.WriteString(" OFFSET $")
		query.WriteString(strconv.Itoa(argIndex + 1))
	}

	args = append(args, limit, offset)

	sqlQuery := query.String()
	ctxlog.L(ctx).Info("listWithTrigramSearch SQL",
		slog.String("query", sqlQuery),
		slog.Any("args", args),
	)

	// Для прямого SQL нужно создать обертку, реализующую SQLBuilder
	// Пока используем старый метод, так как это прямой SQL
	words, err := database.Select[model.Word](ctx, r.q, sqlQuery, args...)
	if err != nil {
		ctxlog.L(ctx).Error("listWithTrigramSearch error",
			slog.String("error", err.Error()),
			slog.String("query", sqlQuery),
			slog.Any("args", args),
		)
		return nil, err
	}

	return words, nil
}

func (r *Repo) Count(ctx context.Context, filter *model.WordFilter) (int, error) {
	// Если есть поиск, используем триграммный поиск с прямой SQL
	if filter != nil && filter.Search != nil && *filter.Search != "" {
		return r.countWithTrigramSearch(ctx, filter)
	}

	// Используем COUNT(DISTINCT words.id) для правильного подсчета при JOIN
	qb := database.Builder.
		Select("COUNT(DISTINCT " + string(schema.Words.ID) + ")").
		From(schema.Words.Name.String())
	qb = applyFilter(qb, filter)

	return database.NewQuery[int](r.q, qb).Scalar(ctx)
}

// countWithTrigramSearch выполняет подсчет с использованием триграммного поиска
func (r *Repo) countWithTrigramSearch(ctx context.Context, filter *model.WordFilter) (int, error) {
	searchQuery := *filter.Search
	similarityThreshold := 0.2

	// Начинаем формировать SQL запрос для COUNT
	var query strings.Builder
	query.WriteString("SELECT COUNT(DISTINCT ")
	query.WriteString(string(schema.Words.ID))
	query.WriteString(") FROM ")
	query.WriteString(schema.Words.Name.String())

	args := []interface{}{searchQuery}
	argIndex := 2 // $1 уже используется для searchQuery

	// Добавляем WHERE условия для триграммного поиска
	whereConditions := []string{
		"word_similarity($1, " + string(schema.Words.Text) + ") > $" + strconv.Itoa(argIndex) + " OR " +
			string(schema.Words.Text) + " % $1 OR " +
			string(schema.Words.Text) + " ILIKE $" + strconv.Itoa(argIndex+1),
	}
	args = append(args, similarityThreshold, "%"+searchQuery+"%")
	argIndex += 2

	// Добавляем JOIN и условия для статуса
	if filter.Status != nil {
		query.WriteString(" JOIN ")
		query.WriteString(schema.Meanings.Name.String())
		query.WriteString(" ON ")
		query.WriteString(string(schema.Meanings.WordID))
		query.WriteString(" = ")
		query.WriteString(string(schema.Words.ID))
		whereConditions = append(whereConditions, string(schema.Meanings.LearningStatus)+" = $"+strconv.Itoa(argIndex))
		args = append(args, *filter.Status)
		argIndex++
	}

	// Добавляем JOIN и условия для тегов
	if len(filter.Tags) > 0 {
		if filter.Status == nil {
			query.WriteString(" JOIN ")
			query.WriteString(schema.Meanings.Name.String())
			query.WriteString(" ON ")
			query.WriteString(string(schema.Meanings.WordID))
			query.WriteString(" = ")
			query.WriteString(string(schema.Words.ID))
		}
		query.WriteString(" JOIN ")
		query.WriteString(schema.MeaningTags.Name.String())
		query.WriteString(" ON ")
		query.WriteString(string(schema.MeaningTags.MeaningID))
		query.WriteString(" = ")
		query.WriteString(string(schema.Meanings.ID))
		query.WriteString(" JOIN ")
		query.WriteString(schema.Tags.Name.String())
		query.WriteString(" ON ")
		query.WriteString(string(schema.Tags.ID))
		query.WriteString(" = ")
		query.WriteString(string(schema.MeaningTags.TagID))

		// Формируем IN условие для тегов
		tagPlaceholders := make([]string, len(filter.Tags))
		for i := range filter.Tags {
			tagPlaceholders[i] = "$" + strconv.Itoa(argIndex)
			args = append(args, filter.Tags[i])
			argIndex++
		}
		whereConditions = append(whereConditions, string(schema.Tags.NameCol)+" IN ("+strings.Join(tagPlaceholders, ", ")+")")
	}

	// Объединяем WHERE условия
	if len(whereConditions) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(strings.Join(whereConditions, " AND "))
	}

	sqlQuery := query.String()
	ctxlog.L(ctx).Info("countWithTrigramSearch SQL",
		slog.String("query", sqlQuery),
		slog.Any("args", args),
	)

	// Для прямого SQL нужно создать обертку, реализующую SQLBuilder
	// Пока используем старый метод, так как это прямой SQL
	count, err := database.GetScalar[int](ctx, r.q, sqlQuery, args...)
	if err != nil {
		ctxlog.L(ctx).Error("countWithTrigramSearch error",
			slog.String("error", err.Error()),
			slog.String("query", sqlQuery),
			slog.Any("args", args),
		)
		return 0, err
	}

	return count, nil
}

func (r *Repo) Exists(ctx context.Context, id int64) (bool, error) {
	builder := database.Builder.
		Select("1").
		From(schema.Words.Name.String()).
		Where(schema.Words.ID.Eq(id)).
		Limit(1)

	val, err := database.NewQuery[int](r.q, builder).Scalar(ctx)
	if err != nil {
		if err == database.ErrNotFound {
			return false, nil
		}
		return false, err
	}
	return val > 0, nil
}

// SearchSimilar использует триграммный поиск для поиска похожих слов.
// Возвращает слова, отсортированные по similarity (от большего к меньшему).
// similarityThreshold - минимальный порог схожести (0.0 - 1.0), по умолчанию 0.3
func (r *Repo) SearchSimilar(ctx context.Context, query string, limit int, similarityThreshold float64) ([]model.Word, error) {
	if limit <= 0 {
		limit = 10
	}
	if similarityThreshold <= 0 {
		similarityThreshold = 0.3 // Порог по умолчанию
	}

	// Используем word_similarity для более точного поиска
	// word_similarity ищет похожие слова, а не подстроки
	// Формируем SELECT с similarity для сортировки
	selectCols := make([]string, 0, len(schema.Words.All()))
	for _, col := range schema.Words.All() {
		selectCols = append(selectCols, string(col))
	}

	// Используем прямой SQL запрос, так как squirrel не поддерживает функции pg_trgm напрямую
	// word_similarity(query, text) > threshold AND text % query (оператор % для триграмм)
	sqlQuery := `
		SELECT ` + strings.Join(selectCols, ", ") + `
		FROM ` + schema.Words.Name.String() + `
		WHERE word_similarity($1, ` + string(schema.Words.Text) + `) > $2
		   OR ` + string(schema.Words.Text) + ` % $1
		ORDER BY word_similarity($1, ` + string(schema.Words.Text) + `) DESC
		LIMIT $3
	`

	args := []interface{}{query, similarityThreshold, limit}

	// Для прямого SQL нужно создать обертку, реализующую SQLBuilder
	// Пока используем старый метод, так как это прямой SQL
	words, err := database.Select[model.Word](ctx, r.q, sqlQuery, args...)
	if err != nil {
		return nil, err
	}

	return words, nil
}

func applyFilter(qb squirrel.SelectBuilder, filter *model.WordFilter) squirrel.SelectBuilder {
	if filter == nil {
		return qb
	}

	// Фильтр по поиску - обрабатывается отдельно в listWithTrigramSearch
	// Здесь оставляем пустым, так как триграммный поиск требует прямой SQL
	if filter.Search != nil && *filter.Search != "" {
		// Пропускаем здесь, обработаем в List через listWithTrigramSearch
	}

	// Фильтр по статусу изучения (требует JOIN с meanings)
	if filter.Status != nil {
		qb = qb.
			Join(schema.Meanings.Name.String() + " ON " + string(schema.Meanings.WordID) + " = " + string(schema.Words.ID)).
			Where(schema.Meanings.LearningStatus.Eq(*filter.Status))
	}

	// Фильтр по тегам (требует JOIN с meanings_tags и tags)
	if len(filter.Tags) > 0 {
		// JOIN с meanings для связи word -> meaning
		needsMeaningsJoin := filter.Status == nil
		if needsMeaningsJoin {
			qb = qb.Join(schema.Meanings.Name.String() + " ON " + string(schema.Meanings.WordID) + " = " + string(schema.Words.ID))
		}

		// JOIN с meanings_tags
		qb = qb.Join(schema.MeaningTags.Name.String() + " ON " + string(schema.MeaningTags.MeaningID) + " = " + string(schema.Meanings.ID))

		// JOIN с tags
		qb = qb.Join(schema.Tags.Name.String() + " ON " + string(schema.Tags.ID) + " = " + string(schema.MeaningTags.TagID))

		// WHERE по именам тегов
		qb = qb.Where(schema.Tags.NameCol.In(filter.Tags))
	}

	return qb
}
