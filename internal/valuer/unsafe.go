package valuer

import (
	"database/sql"
	"github.com/DaHuangQwQ/gweb/internal/errs"
	gsql "github.com/DaHuangQwQ/gweb/model"
	"reflect"
	"unsafe"
)

var _ Creator = NewUnsafeValue

type unsafeValuer struct {
	model   *gsql.Model
	address unsafe.Pointer
}

func NewUnsafeValue(model *gsql.Model, val any) Valuer {
	address := reflect.ValueOf(val).UnsafePointer()

	return unsafeValuer{
		model:   model,
		address: address,
	}
}

func (r unsafeValuer) SetColumns(rows *sql.Rows) error {
	cs, err := rows.Columns()
	if err != nil {
		return err
	}

	var vals []any

	for _, c := range cs {
		// c => column
		fd, ok := r.model.ColumnMap[c]
		if !ok {
			return errs.NewErrUnknownColumn(c)
		}

		fdAddress := unsafe.Pointer(uintptr(r.address) + fd.Offset)

		val := reflect.NewAt(fd.Typ, fdAddress)
		vals = append(vals, val.Interface())
	}

	return rows.Scan(vals...)
}
