package gsql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DaHuangQwQ/gsql/internal/errs"
)

var (
	_ Session = (*Tx)(nil)
	_ Session = (*DB)(nil)
)

type Session interface {
	getCore() core
	queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	execContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type Tx struct {
	db *DB
	tx *sql.Tx
}

func (tx *Tx) getCore() core {
	return tx.db.core
}

func (tx *Tx) queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return tx.tx.QueryContext(ctx, query, args...)
}

func (tx *Tx) execContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return tx.tx.ExecContext(ctx, query, args...)
}

func (tx *Tx) DoTx(ctx context.Context, fn func(ctx context.Context, tx *Tx) error, opts *sql.TxOptions) error {
	tx, err := tx.db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	panicked := true

	defer func() {
		if err != nil || panicked {
			er := tx.Rollback()
			err = errs.NewErrFailedToRollbackTx(err, er, panicked)
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(ctx, tx)

	panicked = false

	return err
}

func (tx *Tx) Commit() error {
	return tx.tx.Commit()
}

func (tx *Tx) Rollback() error {
	return tx.tx.Rollback()
}

func (tx *Tx) RollbackIfNotCommit() error {
	err := tx.tx.Rollback()
	if errors.Is(err, sql.ErrTxDone) {
		err = nil
	}
	return err
}
