package gsql

import "strings"

type Deleter[T any] struct {
	builder

	where   []Predicate
	session Session
}

func NewDeleter[T any](session Session) *Deleter[T] {
	base := session.getCore()
	m, err := base.r.Register(new(T))
	if err != nil {
		return nil
	}
	return &Deleter[T]{
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
		session: session,
	}
}

func (d *Deleter[T]) Build() (*Query, error) {
	d.builder.sb.WriteString("DELETE FROM ")
	d.builder.quote(d.model.TableName)

	if len(d.where) > 0 {
		d.builder.sb.WriteString(" WHERE ")

		p := d.where[0]
		for i := 1; i < len(d.where); i++ {
			p.And(d.where[1])
		}

		er := d.buildExpression(p)
		if er != nil {
			return nil, er
		}
	}

	d.builder.sb.WriteByte(';')

	return &Query{
		SQL:  d.builder.sb.String(),
		Args: d.builder.args,
	}, nil
}

func (d *Deleter[T]) Where(p ...Predicate) *Deleter[T] {
	d.where = p
	return d
}

func (d *Deleter[T]) From(tableName string) *Deleter[T] {
	d.model.TableName = tableName
	return d
}
