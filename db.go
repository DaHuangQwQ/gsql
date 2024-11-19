package gsql

import (
	"database/sql"
	"github.com/DaHuangQwQ/gweb/internal/valuer"
	"github.com/DaHuangQwQ/gweb/model"
)

type DBOption func(db *DB)

type DB struct {
	r       model.Registry
	db      *sql.DB
	creator valuer.Creator
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
		r:       model.NewRegistry(),
		db:      db,
		creator: valuer.NewUnsafeValue,
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

func WithReflectValuer(v valuer.Creator) DBOption {
	return func(db *DB) {
		db.creator = v
	}
}
