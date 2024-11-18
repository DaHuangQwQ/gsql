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

	ErrNoRows = errors.New("no rows in result set")
)

func NewErrUnknownField(name any) error {
	return fmt.Errorf("gsql: unknown field: %v", name)
}

func NewErrUnknownColumn(name any) error {
	return fmt.Errorf("gsql: unknown column: %v", name)
}

func NewErrInvalidTagContent(name any) error {
	return fmt.Errorf("gsql: invalid tag content: %v", name)
}
