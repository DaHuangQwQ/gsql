package gsql

import (
	"context"
	"github.com/DaHuangQwQ/gsql/internal/errs"
	"strings"
)

type Selectable interface {
	selectable()
}

type Selector[T any] struct {
	builder
	table   string
	columns []Selectable
	where   []Predicate

	session Session
}

func NewSelector[T any](db Session) *Selector[T] {
	base := db.getCore()
	m, err := base.r.Register(new(T))
	if err != nil {
		panic(err)
	}

	return &Selector[T]{
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
	}
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	res := get[T](ctx, s.session, s.core, &QueryContext{
		Type:    TypeSelect,
		Builder: s,
		Model:   s.model,
	})

	if res.Result != nil {
		return res.Result.(*T), nil
	}

	return nil, res.Err
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	q, err := s.Build()
	if err != nil {
		return nil, err
	}

	row, err := s.session.queryContext(ctx, q.SQL, q.Args...)
	if err != nil {
		return nil, err
	}

	res := make([]*T, 0, 8)

	for row.Next() {
		tp := new(T)

		val := s.builder.core.creator(s.model, tp)

		err = val.SetColumns(row)
		if err != nil {
			return nil, err
		}

		res = append(res, tp)
	}

	if len(res) == 0 {
		return nil, errs.ErrNoRows
	}

	return res, nil
}

func (s *Selector[T]) Build() (*Query, error) {
	s.sb.WriteString("SELECT ")

	if err := s.buildColumns(); err != nil {
		return nil, err
	}

	s.sb.WriteString(" FROM ")

	if s.table == "" {
		s.From(s.model.TableName)
	}
	s.sb.WriteString(s.table)

	if len(s.where) > 0 {
		s.sb.WriteString(" WHERE ")

		p := s.where[0]
		for i := 1; i < len(s.where); i++ {
			p = p.And(s.where[i])
		}

		er := s.buildExpression(p)
		if er != nil {
			return nil, er
		}
	}

	s.sb.WriteByte(';')

	return &Query{
		SQL:  s.sb.String(),
		Args: s.args,
	}, nil
}

func (s *Selector[T]) buildColumns() error {
	if len(s.columns) == 0 {
		s.sb.WriteByte('*')
		return nil
	}

	for i, col := range s.columns {
		if i > 0 {
			s.sb.WriteByte(',')
		}
		switch c := col.(type) {
		case Column:
			er := s.buildColumn(c)
			if er != nil {
				return er
			}
		case Aggregate:
			s.sb.WriteString(c.fn)
			s.sb.WriteByte('(')
			er := s.buildColumn(Column{
				Name: c.arg,
			})
			s.sb.WriteByte(')')
			if er != nil {
				return er
			}
			// 聚合函数的别名
			if c.alias != "" {
				s.sb.WriteString(" AS ")
				s.quote(c.alias)
			}
		case RawExpr:
			s.sb.WriteString(c.raw)
			s.addArgs(c.args...)
		}
	}

	return nil
}

func (s *Selector[T]) From(table string) *Selector[T] {
	if table == "" {
		return s
	}

	str := ""
	segs := strings.Split(table, ".")

	for i := 0; i < len(segs); i++ {
		if i < len(segs)-1 {
			str += "`" + segs[i] + "`."
		} else {
			str += "`" + segs[i] + "`"
		}
	}
	s.table = str
	return s
}

func (s *Selector[T]) Select(cols ...Selectable) *Selector[T] {
	s.columns = cols
	return s
}

func (s *Selector[T]) Where(p ...Predicate) *Selector[T] {
	s.where = p
	return s
}
