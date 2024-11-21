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

	ErrInsertZeroRow = errors.New("no values to insert")
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

func NewErrUnsupportedExpression(expr any) error {
	return fmt.Errorf("gsql: unsupported expression: %v", expr)
}

func NewErrUnsupportedTable(table any) error {
	return fmt.Errorf("gsql: unsupported table: %v", table)
}

func NewErrUnsupportedAssignable(assign any) error {
	return fmt.Errorf("gsql: unsupported assignable: %v", assign)
}
