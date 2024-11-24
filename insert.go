package gsql

import (
	"context"
	"github.com/DaHuangQwQ/gsql/internal/errs"
	"github.com/DaHuangQwQ/gsql/model"
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

	session        Session
	onDuplicateKey *Upsert
}

func NewInserter[T any](db Session) *Inserter[T] {
	base := db.getCore()
	m, err := base.r.Register(new(T))
	if err != nil {
		panic(err)
	}

	return &Inserter[T]{
		builder: builder{
			core: core{
				model:   m,
				dialect: base.dialect,
				creator: base.creator,
				r:       base.r,
				mdls:    base.mdls,
			},
			sb:     strings.Builder{},
			quoter: base.dialect.quoter(),
		},
		session: db,
		values:  []*T{},
	}
}

func (i *Inserter[T]) Build() (*Query, error) {
	if len(i.values) == 0 {
		return nil, errs.ErrInsertZeroRow
	}

	i.sb.WriteString("INSERT INTO ")

	i.quote(i.core.model.TableName)

	fields := i.core.model.Fields

	if len(i.columns) > 0 {
		fields = make([]*model.Field, 0, len(i.columns))
		for _, col := range i.columns {
			val, ok := i.core.model.FieldMap[col]
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

func (i *Inserter[T]) Exec(ctx context.Context) Result {

	res := exec(ctx, i.session, i.core, &QueryContext{
		Type:    TypeInsert,
		Builder: i,
		Model:   i.model,
	})

	if res.Result != nil {
		return Result{
			res: res.Result.(Result),
		}
	}

	return Result{
		err: res.Err,
	}
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
