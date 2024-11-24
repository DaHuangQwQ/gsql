package gsql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/DaHuangQwQ/gsql/internal/valuer"
	"github.com/DaHuangQwQ/gsql/model"
	"log"
)

type DBOption func(db *DB)

type DB struct {
	core
	db *sql.DB
}

func (db *DB) Use(mdls ...Middleware) {
	db.core.mdls = append(db.core.mdls, mdls...)
}

func (db *DB) getCore() core {
	return db.core
}

func (db *DB) queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return db.db.QueryContext(ctx, query, args...)
}

func (db *DB) execContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return db.db.ExecContext(ctx, query, args...)
}

func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &Tx{
		tx: tx,
	}, nil
}

func (db *DB) Wait() error {
	err := db.db.Ping()
	for errors.Is(err, driver.ErrBadConn) {
		log.Println("gsql: err bad connection")
		err = db.db.Ping()
	}
	return nil
}

func Open(driver string, dsn string, opts ...DBOption) (*DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	return OpenDB(db, opts...)
}

func OpenDB(db *sql.DB, opts ...DBOption) (*DB, error) {
	res := &DB{
		core: core{
			r:       model.NewRegistry(),
			creator: valuer.NewUnsafeValue,
			dialect: DialectMySQL,
		},
		db: db,
	}

	for _, opt := range opts {
		opt(res)
	}

	return res, nil
}

func MustOpen(driver string, dsn string, opts ...DBOption) *DB {
	db, err := Open(driver, dsn, opts...)
	if err != nil {
		panic(err)
	}
	return db
}

func WithValuer(v valuer.Creator) DBOption {
	return func(db *DB) {
		db.creator = v
	}
}

func WithRegistry(r model.Registry) DBOption {
	return func(db *DB) {
		db.r = r
	}
}

func WithDialect(dialect Dialect) DBOption {
	return func(db *DB) {
		db.dialect = dialect
	}
}
