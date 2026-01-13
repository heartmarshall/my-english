package transport

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
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
	return nil
}
