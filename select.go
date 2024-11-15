package gsql

import (
	"context"
	"github.com/DaHuangQwQ/gweb/internal/errs"
	"strings"
)

type Selector[T any] struct {
	table string
	model *model
	where []Predicate
	sb    strings.Builder
	args  []any

	db *DB
}

func NewSelector[T any](db *DB) *Selector[T] {
	return &Selector[T]{
		sb: strings.Builder{},
		db: db,
	}
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Selector[T]) Build() (*Query, error) {
	var err error
	s.model, err = s.db.r.ParseModel(new(T))
	if err != nil {
		return nil, err
	}

	s.sb.WriteString("SELECT * FROM ")

	if s.table == "" {
		s.From(s.model.tableName)
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

func (s *Selector[T]) addArgs(val any) {
	if s.args == nil {
		s.args = make([]any, 0, 4)
	}
	s.args = append(s.args, val)
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
		s.sb.WriteByte(' ')
		s.sb.WriteString(exp.op.String())
		s.sb.WriteByte(' ')

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
		fd, ok := s.model.fields[exp.Name]
		if !ok {
			return errs.NewErrUnknownField(exp.Name)
		}
		s.sb.WriteByte('`')
		s.sb.WriteString(fd.colName)
		s.sb.WriteByte('`')
		return nil
	case Value:
		s.sb.WriteByte('?')
		s.addArgs(exp.val)
		return nil
	default:
		return errs.ErrInvalidExpression
	}
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

func (s *Selector[T]) Where(p ...Predicate) *Selector[T] {
	s.where = p
	return s
}
