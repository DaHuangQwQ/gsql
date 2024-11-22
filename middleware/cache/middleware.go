package cache

import (
	"context"
	"github.com/DaHuangQwQ/gsql"
)

type MiddlewareBuilder struct {
}

func (m *MiddlewareBuilder) Build() gsql.Middleware {
	return func(next gsql.Handler) gsql.Handler {
		return func(ctx context.Context, qc *gsql.QueryContext) *gsql.QueryResult {
			//if qc.Type != "SELECT" {
			//	return next(ctx, qc)
			//}
			//bd := qc.Builder.(gsql.Selector[User])
			//tr := bd.Table().(gsql.Table)
			//
			//// 从缓存读到了
			//if readFromCache() {
			//	return &gsql.QueryResult{}
			//}
			//
			return next(ctx, qc)
		}
	}
}
