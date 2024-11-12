package gsql

import (
	"context"
	"reflect"
	"strings"
)

type Selector[T any] struct {
	table string
	where []Predicate
	sb    strings.Builder
	args  []any
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
	s.sb.WriteString("SELECT * FROM ")

	if s.table == "" {
		var t T
		tye := reflect.TypeOf(t)
		s.From(tye.Name())
	}
	s.sb.WriteString(s.table)

	if len(s.where) > 0 {
		s.sb.WriteString(" WHERE ")

		p := s.where[0]
		for i := 1; i < len(s.where); i++ {
			p = p.And(s.where[i])
		}

		s.buildExpression(p)
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

func (s *Selector[T]) buildExpression(expr Expression) {
	switch exp := expr.(type) {
	case nil:
		return
	case Predicate:
		_, ok := exp.left.(Predicate)
		if ok {
			s.sb.WriteByte('(')
		}
		s.buildExpression(exp.left)
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
		s.buildExpression(exp.right)
		if ok {
			s.sb.WriteByte(')')
		}
	case Column:
		s.sb.WriteByte('`')
		s.sb.WriteString(exp.Name)
		s.sb.WriteByte('`')
	case Value:
		s.sb.WriteByte('?')
		s.addArgs(exp.val)
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
