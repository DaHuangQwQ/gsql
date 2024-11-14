package gsql

type DBOption func(db *DB)

type DB struct {
	r *registry
}

func NewDB(opts ...DBOption) (*DB, error) {
	res := &DB{
		r: newRegistry(),
	}

	for _, opt := range opts {
		opt(res)
	}

	return res, nil
}

func MustNewDB(opts ...DBOption) *DB {
	db, err := NewDB(opts...)
	if err != nil {
		panic(err)
	}
	return db
}
