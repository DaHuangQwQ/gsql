package gsql

import (
	"github.com/DaHuangQwQ/gweb/internal/errs"
	"github.com/DaHuangQwQ/gweb/model"
	"strings"
)

type UpsertBuilder[T any] struct {
	i               *Inserter[T]
	conflictColumns []string
}

func (o *UpsertBuilder[T]) ConflictColumns(cols ...string) *UpsertBuilder[T] {
	o.conflictColumns = cols
	return o
}

func (o *UpsertBuilder[T]) Update(assigns ...Assignable) *Inserter[T] {
	o.i.onDuplicateKey = &Upsert{
		assigns:         assigns,
		conflictColumns: o.conflictColumns,
	}
	return o.i
}

type Upsert struct {
	assigns         []Assignable
	conflictColumns []string
}

type Inserter[T any] struct {
	builder
	values  []*T
	columns []string
	model   *model.Model

	db             *DB
	onDuplicateKey *Upsert
}

func NewInserter[T any](db *DB) *Inserter[T] {
	m, err := db.r.Register(new(T))
	if err != nil {
		panic(err)
	}

	return &Inserter[T]{
		builder: builder{
			core: core{
				model:   m,
				dialect: db.dialect,
				creator: db.creator,
				r:       db.r,
			},
			sb:     strings.Builder{},
			quoter: db.dialect.quoter(),
		},
		db:     db,
		values: []*T{},
	}
}

func (i *Inserter[T]) Build() (*Query, error) {
	if len(i.values) == 0 {
		return nil, errs.ErrInsertZeroRow
	}

	if i.model == nil {
		m, err := i.db.r.Get(i.values[0])
		if err != nil {
			return nil, err
		}
		i.model = m
	}

	i.sb.WriteString("INSERT INTO ")

	i.quote(i.model.TableName)

	fields := i.model.Fields

	if len(i.columns) > 0 {
		fields = make([]*model.Field, 0, len(i.columns))
		for _, col := range i.columns {
			val, ok := i.model.FieldMap[col]
			if !ok {
				return nil, errs.NewErrUnknownField(col)
			}
			fields = append(fields, val)
		}
	}

	i.sb.WriteByte('(')
	for idx, val := range fields {
		if idx > 0 {
			i.sb.WriteByte(',')
		}
		i.quote(val.ColName)
	}
	i.sb.WriteByte(')')

	i.sb.WriteString(" VALUES ")

	i.args = make([]any, 0, len(i.values))

	for idx, value := range i.values {
		if idx > 0 {
			i.sb.WriteByte(',')
		}

		i.sb.WriteByte('(')

		val := i.creator(i.model, value)

		for j, field := range fields {
			if j > 0 {
				i.sb.WriteByte(',')
			}
			i.sb.WriteByte('?')

			arg, err := val.Field(field.GoName)
			if err != nil {
				return nil, errs.NewErrUnknownField(field.GoName)
			}

			i.addArgs(arg)
		}
		i.sb.WriteByte(')')
	}

	if i.onDuplicateKey != nil {
		err := i.dialect.buildUpsert(&i.builder, i.onDuplicateKey)
		if err != nil {
			return nil, err
		}
	}

	i.sb.WriteByte(';')

	return &Query{
		SQL:  i.sb.String(),
		Args: i.args,
	}, nil
}

func (i *Inserter[T]) OnDuplicateKey() *UpsertBuilder[T] {
	return &UpsertBuilder[T]{
		i: i,
	}
}

func (i *Inserter[T]) Columns(columns ...string) *Inserter[T] {
	i.columns = columns
	return i
}

func (i *Inserter[T]) Values(values ...*T) *Inserter[T] {
	i.values = values
	return i
}
