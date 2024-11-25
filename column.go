package gsql

type Column struct {
	Name  string
	alias string
	Table TableReference
}

func C(name string) Column {
	return Column{Name: name}
}

func (c Column) Eq(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opEQ,
		right: valueOf(arg),
	}
}

func (c Column) As(alias string) Column {
	return Column{
		Name:  c.Name,
		alias: alias,
		Table: c.Table,
	}
}

func valueOf(arg any) Expression {
	switch val := arg.(type) {
	case Expression:
		return val
	default:
		return Value{val: val}
	}
}

func (c Column) expr() {}

func (c Column) selectable() {}

func (c Column) assign() {}
