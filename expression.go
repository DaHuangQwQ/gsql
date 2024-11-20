package gsql

// Expression 标记接口， 代表表达式
type Expression interface {
	expr()
}

// RawExpr 原生表达式
type RawExpr struct {
	raw  string
	args []any
}

func Raw(expr string, args ...any) RawExpr {
	return RawExpr{
		raw:  expr,
		args: args,
	}
}

func (r RawExpr) selectable() {}
func (r RawExpr) expr()       {}

func (r RawExpr) AsPredicate() Predicate {
	return Predicate{
		left: r,
	}
}

type binaryExpr struct {
	left  Expression
	op    op
	right Expression
}

func (binaryExpr) expr() {}

type MathExpr binaryExpr

//func (m MathExpr) Add(val interface{}) MathExpr {
//	return MathExpr{
//		left:  m,
//		op:    opAdd,
//		right: valueOf(val),
//	}
//}
//
//func (m MathExpr) Multi(val interface{}) MathExpr {
//	return MathExpr{
//		left:  m,
//		op:    opMulti,
//		right: valueOf(val),
//	}
//}

func (m MathExpr) expr() {}
