package types

import "errors"

var (
	ErrNotFound      = errors.New("entity not found")
	ErrAlreadyExists = errors.New("entity already exists")
	ErrInvalidInput  = errors.New("invalid input data")
	ErrInternal      = errors.New("internal error")
)
