package gsql

import (
	"github.com/DaHuangQwQ/gsql/internal/errs"
	"strings"
)

type builder struct {
	core
	sb     strings.Builder
	args   []any
	quoter byte
}

func (b *builder) quote(name string) {
	b.sb.WriteByte(b.quoter)
	b.sb.WriteString(name)
	b.sb.WriteByte(b.quoter)
}

func (b *builder) buildColumn(col Column) error {
	switch table := col.Table.(type) {
	case nil:
		fd, ok := b.model.FieldMap[col.Name]
		if !ok {
			return errs.NewErrUnknownField(col.Name)
		}
		b.quote(fd.ColName)
		if col.alias != "" {
			b.sb.WriteString(" AS ")
			b.quote(col.alias)
		}
		return nil
	case Table:
		m, err := b.r.Get(table.entity)
		if err != nil {
			return err
		}
		fd, ok := m.FieldMap[col.Name]
		if !ok {
			return errs.NewErrUnknownField(col.Name)
		}
		if table.alias != "" {
			b.quote(table.alias)
			b.sb.WriteByte('.')
		}
		b.quote(fd.ColName)
		if col.alias != "" {
			b.sb.WriteString(" AS ")
			b.quote(col.alias)
		}
		return nil
	default:
		return errs.NewErrUnsupportedTable(col.Name)
	}
}

func (b *builder) addArgs(vals ...any) {
	if len(vals) == 0 {
		return
	}
	if b.args == nil {
		b.args = make([]any, 0, 8)
	}
	b.args = append(b.args, vals...)
}

func (b *builder) buildPredicates(ps []Predicate) error {
	p := ps[0]
	for i := 1; i < len(ps); i++ {
		p = p.And(ps[i])
	}
	return b.buildExpression(p)
}

func (b *builder) buildExpression(expr Expression) error {
	switch exp := expr.(type) {
	case nil:
		return nil
	case Predicate:
		_, ok := exp.left.(Predicate)
		if ok {
			b.sb.WriteByte('(')
		}
		err := b.buildExpression(exp.left)
		if err != nil {
			return err
		}
		if ok {
			b.sb.WriteByte(')')
		}

		if exp.op != "" {
			b.sb.WriteByte(' ')
			b.sb.WriteString(exp.op.String())
			b.sb.WriteByte(' ')
		}

		_, ok = exp.right.(Predicate)
		if ok {
			b.sb.WriteByte('(')
		}
		err = b.buildExpression(exp.right)
		if err != nil {
			return err
		}
		if ok {
			b.sb.WriteByte(')')
		}
		return nil
	case Column:
		exp.alias = ""
		return b.buildColumn(exp)
	case Value:
		b.sb.WriteByte('?')
		b.addArgs(exp.val)
		return nil
	case RawExpr:
		b.sb.WriteByte('(')
		b.sb.WriteString(exp.raw)
		b.addArgs(exp.args...)
		b.sb.WriteByte(')')
		return nil
	default:
		return errs.ErrInvalidExpression
	}
}

//func (b *builder) buildBinaryExpr(e binaryExpr) error {
//	err := b.buildSubExpr(e.left)
//	if err != nil {
//		return err
//	}
//	if e.op != "" {
//		b.sb.WriteByte(' ')
//		b.sb.WriteString(e.op.String())
//	}
//	if e.right != nil {
//		b.sb.WriteByte(' ')
//		return b.buildSubExpr(e.right)
//	}
//	return nil
//}

//func (b *builder) buildSubExpr(subExpr Expression) error {
//	switch sub := subExpr.(type) {
//	case MathExpr:
//		_ = b.sb.WriteByte('(')
//		if err := b.buildBinaryExpr(binaryExpr(sub)); err != nil {
//			return err
//		}
//		_ = b.sb.WriteByte(')')
//	case binaryExpr:
//		_ = b.sb.WriteByte('(')
//		if err := b.buildBinaryExpr(sub); err != nil {
//			return err
//		}
//		_ = b.sb.WriteByte(')')
//	case Predicate:
//		_ = b.sb.WriteByte('(')
//		if err := b.buildBinaryExpr(binaryExpr(sub)); err != nil {
//			return err
//		}
//		_ = b.sb.WriteByte(')')
//	default:
//		if err := b.buildExpression(sub); err != nil {
//			return err
//		}
//	}
//	return nil
//}

func (b *builder) reset() {
	b.sb.Reset()
	b.args = nil
}
