package gsql

import (
	"context"
	"github.com/DaHuangQwQ/gweb/internal/errs"
	"reflect"
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

	cs, err := row.Columns()
	if err != nil {
		return nil, err
	}

	tp := new(T)

	vals := make([]any, 0, len(cs))
	valElems := make([]reflect.Value, 0, len(cs))

	for _, c := range cs {
		// c 是列名
		fd, ok := s.model.columnMap[c]
		if !ok {
			return nil, errs.NewErrUnknownColumn(c)
		}

		// 反射创建一个实例 这里创建的实例是原本类型的指针类型
		val := reflect.New(fd.typ)
		vals = append(vals, val.Interface())
		valElems = append(valElems, val.Elem())
	}

	// 给 vals 里的赋值
	// 类型要匹配 顺序要匹配
	err = row.Scan(vals...)
	if err != nil {
		return nil, err
	}

	tpValueElem := reflect.ValueOf(tp).Elem()
	for i, c := range cs {
		// c 是列名
		fd, ok := s.model.columnMap[c]
		if !ok {
			return nil, errs.NewErrUnknownColumn(c)
		}
		tpValueElem.FieldByName(fd.goName).
			Set(valElems[i])
	}

	return tp, nil
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

		cs, er := row.Columns()
		if er != nil {
			return nil, er
		}

		vals := make([]any, 0, len(cs))
		valElems := make([]reflect.Value, 0, len(cs))

		for _, c := range cs {
			// c 是列名
			fd, ok := s.model.columnMap[c]
			if !ok {
				return nil, errs.NewErrUnknownColumn(c)
			}

			// 反射创建一个实例 这里创建的实例是原本类型的指针类型
			val := reflect.New(fd.typ)
			vals = append(vals, val.Interface())
			valElems = append(valElems, val.Elem())
		}

		// 给 vals 里的赋值
		// 类型要匹配 顺序要匹配
		er = row.Scan(vals...)
		if er != nil {
			return nil, er
		}

		tpValueElem := reflect.ValueOf(tp).Elem()
		for i, c := range cs {
			// c 是列名
			fd, ok := s.model.columnMap[c]
			if !ok {
				return nil, errs.NewErrUnknownColumn(c)
			}
			tpValueElem.FieldByName(fd.goName).
				Set(valElems[i])
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
		fd, ok := s.model.fieldMap[exp.Name]
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
