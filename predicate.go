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
