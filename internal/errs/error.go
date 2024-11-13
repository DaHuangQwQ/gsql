package errs

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidType model type is invalid
	ErrInvalidType = errors.New("invalid type")

	// ErrInvalidExpression sql expression is invalid after where
	ErrInvalidExpression = errors.New("invalid expression")
)

func NewErrUnknownField(name any) error {
	return fmt.Errorf("gsql: unknown field: %v", name)
}
