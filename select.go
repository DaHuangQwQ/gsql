package gsql

import (
	"context"
	"github.com/DaHuangQwQ/gweb/internal/errs"
	model2 "github.com/DaHuangQwQ/gweb/model"
	"strings"
)

type Selectable interface {
	selectable()
}

type Selector[T any] struct {
	table   string
	model   *model2.Model
	columns []Selectable
	where   []Predicate
	sb      strings.Builder
	args    []any

	db *DB
}

func NewSelector[T any](db *DB) *Selector[T] {
	return &Selector[T]{
		sb: strings.Builder{},
		db: db,
	}
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	q, err := s.Build()
	if err != nil {
		return nil, err
	}

	db := s.db.db

	row, err := db.QueryContext(ctx, q.SQL, q.Args...)
	if err != nil {
		return nil, err
	}

	if !row.Next() {
		return nil, errs.ErrNoRows
	}

	tp := new(T)

	val := s.db.creator(s.model, tp)

	err = val.SetColumns(row)

	return tp, err
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	q, err := s.Build()
	if err != nil {
		return nil, err
	}

	db := s.db.db

	row, err := db.QueryContext(ctx, q.SQL, q.Args...)
	if err != nil {
		return nil, err
	}

	res := make([]*T, 0, 8)

	for row.Next() {
		tp := new(T)

		val := s.db.creator(s.model, tp)

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
	var err error
	s.model, err = s.db.r.Register(new(T))
	if err != nil {
		return nil, err
	}

	s.sb.WriteString("SELECT ")

	if err = s.buildColumns(); err != nil {
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

func (s *Selector[T]) addArgs(vals ...any) {
	if len(vals) == 0 {
		return
	}
	if s.args == nil {
		s.args = make([]any, 0, 4)
	}
	s.args = append(s.args, vals...)
}

func (s *Selector[T]) buildExpression(expr Expression) error {
	switch exp := expr.(type) {
	case nil:
		return nil
	case Predicate:
		_, ok := exp.left.(Predicate)
		if ok {
			s.sb.WriteByte('(')
		}
		err := s.buildExpression(exp.left)
		if err != nil {
			return err
		}
		if ok {
			s.sb.WriteByte(')')
		}

		if exp.op != "" {
			s.sb.WriteByte(' ')
			s.sb.WriteString(exp.op.String())
			s.sb.WriteByte(' ')
		}

		_, ok = exp.right.(Predicate)
		if ok {
			s.sb.WriteByte('(')
		}
		err = s.buildExpression(exp.right)
		if err != nil {
			return err
		}
		if ok {
			s.sb.WriteByte(')')
		}
		return nil
	case Column:
		exp.alias = ""
		return s.buildColumn(exp)
	case Value:
		s.sb.WriteByte('?')
		s.addArgs(exp.val)
		return nil
	case RawExpr:
		s.sb.WriteByte('(')
		s.sb.WriteString(exp.raw)
		s.addArgs(exp.args...)
		s.sb.WriteByte(')')
		return nil
	default:
		return errs.ErrInvalidExpression
	}
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
				s.sb.WriteByte('`')
				s.sb.WriteString(c.alias)
				s.sb.WriteByte('`')
			}
		case RawExpr:
			s.sb.WriteString(c.raw)
			s.addArgs(c.args...)
		}
	}

	return nil
}

func (s *Selector[T]) buildColumn(col Column) error {
	fd, ok := s.model.FieldMap[col.Name]
	if !ok {
		return errs.NewErrUnknownField(col.Name)
	}
	s.sb.WriteByte('`')
	s.sb.WriteString(fd.ColName)
	s.sb.WriteByte('`')
	if col.alias != "" {
		s.sb.WriteString(" AS ")
		s.sb.WriteByte('`')
		s.sb.WriteString(col.alias)
		s.sb.WriteByte('`')
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
