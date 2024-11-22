package gsql

import (
	"context"
	"github.com/DaHuangQwQ/gsql/model"
)

type Type string

func (t Type) String() string {
	return string(t)
}

const (
	TypeSelect Type = "select"
	TypeInsert Type = "insert"
	TypeUpdate Type = "update"
	TypeDelete Type = "delete"
)

type QueryContext struct {
	// 查询类型，标记增删改查
	Type Type

	// 代表的是查询本身
	Builder QueryBuilder

	query *Query

	Model *model.Model
}

func (qc *QueryContext) BuildQuery() (*Query, error) {
	var err error
	if qc.query == nil {
		qc.query, err = qc.Builder.Build()
	}
	return qc.query, err
}

type QueryResult struct {
	// Result 在不同查询下类型是不同的
	Result any
	Err    error
}

type Handler func(ctx context.Context, qc *QueryContext) *QueryResult

type Middleware func(next Handler) Handler
