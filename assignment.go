package gsql

type Assignable interface {
	assign()
}

type Assignment struct {
	col string

	// 在 UPDATE 里面改成了这个 Expression 的结构
	val Expression
}

func (Assignment) assign() {}

func Assign(col string, val any) Assignment {
	return Assignment{
		col: col,
		val: valueOf(val),
	}
}
