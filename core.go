package gsql

import (
	"github.com/DaHuangQwQ/gsql/internal/valuer"
	"github.com/DaHuangQwQ/gsql/model"
)

type core struct {
	model   *model.Model
	dialect Dialect
	creator valuer.Creator
	r       model.Registry
	mdls    []Middleware
}
