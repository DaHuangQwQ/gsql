package nodelete

import (
	"context"
	"errors"
	gsql "github.com/DaHuangQwQ/gsql"
)

type MiddlewareBuilder struct {
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{}
}

func (m MiddlewareBuilder) Build() gsql.Middleware {
	return func(next gsql.Handler) gsql.Handler {
		return func(ctx context.Context, qc *gsql.QueryContext) *gsql.QueryResult {
			// 禁用 DELETE 语句
			if qc.Type == "DELETE" {
				return &gsql.QueryResult{
					Err: errors.New("no Delete"),
				}
			}
			return next(ctx, qc)
		}
	}
}
