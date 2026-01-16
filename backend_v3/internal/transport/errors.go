package transport

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/heartmarshall/my-english/internal/service/types"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// ErrorCode — коды ошибок для GraphQL.
type ErrorCode string

const (
	CodeNotFound      ErrorCode = "NOT_FOUND"
	CodeAlreadyExists ErrorCode = "ALREADY_EXISTS"
	CodeInvalidInput  ErrorCode = "INVALID_INPUT"
	CodeInternal      ErrorCode = "INTERNAL_ERROR"
)

// HandleError преобразует service ошибки в GraphQL ошибки.
func HandleError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	// 1. Ошибка валидации
	if types.IsValidationError(err) {
		return &gqlerror.Error{
			Message: err.Error(),
			Path:    graphql.GetPath(ctx),
			Extensions: map[string]interface{}{
				"code": CodeInvalidInput,
			},
		}
	}

	// 2. Not Found
	if errors.Is(err, types.ErrNotFound) {
		return &gqlerror.Error{
			Message: "Entity not found",
			Path:    graphql.GetPath(ctx),
			Extensions: map[string]interface{}{
				"code": CodeNotFound,
			},
		}
	}

	// 3. Already Exists
	if errors.Is(err, types.ErrAlreadyExists) {
		return &gqlerror.Error{
			Message: "Entity already exists",
			Path:    graphql.GetPath(ctx),
			Extensions: map[string]interface{}{
				"code": CodeAlreadyExists,
			},
		}
	}

	// 4. Invalid Input (общий)
	if errors.Is(err, types.ErrInvalidInput) {
		return &gqlerror.Error{
			Message: "Invalid input data",
			Path:    graphql.GetPath(ctx),
			Extensions: map[string]interface{}{
				"code": CodeInvalidInput,
			},
		}
	}

	// 5. Internal Error (скрываем детали, если это не dev mode, но для простоты пока возвращаем как есть)
	// В продакшене лучше логировать err и возвращать "Internal Server Error"
	return &gqlerror.Error{
		Message: "Internal server error",
		Path:    graphql.GetPath(ctx),
		Extensions: map[string]interface{}{
			"code": CodeInternal,
		},
	}
}
