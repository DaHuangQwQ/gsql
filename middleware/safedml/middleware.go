package safedml

import (
	"context"
	"errors"
	"github.com/DaHuangQwQ/gsql"
	"strings"
)

type MiddlewareBuilder struct {
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{}
}

func (m MiddlewareBuilder) Build() gsql.Middleware {
	return func(next gsql.Handler) gsql.Handler {
		return func(ctx context.Context, qc *gsql.QueryContext) *gsql.QueryResult {
			if qc.Type == "SELECT" || qc.Type == "INSERT" {
				return next(ctx, qc)
			}
			q, err := qc.Builder.Build()
			if err != nil {
				return &gsql.QueryResult{
					Err: err,
				}
			}
			if strings.Contains(q.SQL, "WHERE") {
				return &gsql.QueryResult{
					Err: errors.New("不准执行没有 WHERE 的 delete 或者 update 语句"),
				}
			}
			return next(ctx, qc)
		}
	}
}
