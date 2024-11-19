package valuer

import (
	"database/sql"
	gsql "github.com/DaHuangQwQ/gweb/model"
)

type Valuer interface {
	SetColumns(rows *sql.Rows) error
}

type Creator func(model *gsql.Model, entity any) Valuer
