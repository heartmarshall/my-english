package transport

import (
	"context"
	"errors"
	"log/slog"

	"github.com/99designs/gqlgen/graphql"
	"github.com/heartmarshall/my-english/internal/service"
	ctxlog "github.com/heartmarshall/my-english/pkg/context"
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

// NewGraphQLError создаёт GraphQL ошибку с расширениями.
func NewGraphQLError(ctx context.Context, message string, code ErrorCode) *gqlerror.Error {
	return &gqlerror.Error{
		Message: message,
		Path:    graphql.GetPath(ctx),
		Extensions: map[string]interface{}{
			"code": string(code),
		},
	}
}

// HandleError преобразует service ошибки в GraphQL ошибки.
func HandleError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, service.ErrWordNotFound):
		return NewGraphQLError(ctx, "Word not found", CodeNotFound)

	case errors.Is(err, service.ErrMeaningNotFound):
		return NewGraphQLError(ctx, "Meaning not found", CodeNotFound)

	case errors.Is(err, service.ErrWordAlreadyExists):
		return NewGraphQLError(ctx, "Word already exists", CodeAlreadyExists)

	case errors.Is(err, service.ErrInvalidInput):
		return NewGraphQLError(ctx, "Invalid input", CodeInvalidInput)

	case errors.Is(err, service.ErrInvalidGrade):
		return NewGraphQLError(ctx, "Grade must be between 1 and 5", CodeInvalidInput)

	default:
		// Логируем реальную ошибку, но не показываем пользователю
		logger := ctxlog.L(ctx)
		logger.Error("internal server error",
			slog.String("error", err.Error()),
			slog.Any("error_type", err),
		)
		return NewGraphQLError(ctx, "Internal server error", CodeInternal)
	}
}
