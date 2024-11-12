/*
	package gsql

表达式生成
*/
package gsql

type op string

func (o op) String() string {
	return string(o)
}

const (
	opEQ  op = "="
	opLT  op = "<"
	opGT  op = ">"
	opNOT op = "NOT"
	opAND op = "AND"
	opOR  op = "OR"
)

type Predicate struct {
	left  Expression
	right Expression
	op    op
}

type Column struct {
	Name string
}

func C(name string) Column {
	return Column{Name: name}
}

func (c Column) Eq(arg any) Predicate {
	return Predicate{
		left: c,
		op:   opEQ,
		right: Value{
			val: arg,
		},
	}
}

func (c Column) expr() {}

func Not(p Predicate) Predicate {
	return Predicate{
		op:    opNOT,
		right: p,
	}
}

func (left Predicate) And(right Predicate) Predicate {
	return Predicate{
		left:  left,
		op:    opAND,
		right: right,
	}
}

func (left Predicate) Or(right Predicate) Predicate {
	return Predicate{
		left:  left,
		op:    opOR,
		right: right,
	}
}

func (left Predicate) expr() {}

type Value struct {
	val any
}

func (Value) expr() {}

// Expression 标记接口， 代表表达式
type Expression interface {
	expr()
}
