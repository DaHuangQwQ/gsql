package gsql

import (
	"github.com/DaHuangQwQ/gweb/internal/valuer"
	"github.com/DaHuangQwQ/gweb/model"
)

type core struct {
	model   *model.Model
	dialect Dialect
	creator valuer.Creator
	r       model.Registry
}
