package valuer

import (
	"database/sql"
	"github.com/DaHuangQwQ/gweb/internal/errs"
	gsql "github.com/DaHuangQwQ/gweb/model"
	"reflect"
)

var _ Creator = NewReflectValue

type reflectValuer struct {
	model *gsql.Model
	val   reflect.Value
}

func NewReflectValue(model *gsql.Model, val any) Valuer {
	return reflectValuer{
		model: model,
		val:   reflect.ValueOf(val).Elem(),
	}
}

func (r reflectValuer) SetColumns(rows *sql.Rows) error {
	cs, err := rows.Columns()
	if err != nil {
		return err
	}

	vals := make([]any, 0, len(cs))
	valElems := make([]reflect.Value, 0, len(cs))
	for _, c := range cs {
		// c is column
		fd, ok := r.model.ColumnMap[c]
		if !ok {
			return errs.NewErrUnknownColumn(c)
		}

		val := reflect.New(fd.Typ)
		vals = append(vals, val.Interface())
		valElems = append(valElems, val.Elem())
	}

	err = rows.Scan(vals...)
	if err != nil {
		return err
	}

	tpValueElem := r.val
	for i, c := range cs {
		// c is column
		fd, ok := r.model.ColumnMap[c]
		if !ok {
			return errs.NewErrUnknownColumn(c)
		}
		tpValueElem.FieldByName(fd.GoName).
			Set(valElems[i])
	}

	return nil
}
