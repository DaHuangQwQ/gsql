package querylog

import (
	"context"
	gsql "github.com/DaHuangQwQ/gsql"
	"log"
)

type MiddlewareBuilder struct {
	logFunc func(sql string, args []any)
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		logFunc: func(sql string, args []any) {
			log.Printf("gsql: query: %s, args: %v \n", sql, args)
		},
	}
}

func (m *MiddlewareBuilder) Build() gsql.Middleware {
	return func(next gsql.Handler) gsql.Handler {
		return func(ctx context.Context, qc *gsql.QueryContext) *gsql.QueryResult {

			q, err := qc.Builder.Build()
			if err != nil {
				return &gsql.QueryResult{
					Err: err,
				}
			}

			m.logFunc(q.SQL, q.Args)

			return next(ctx, qc)
		}
	}
}

func (m *MiddlewareBuilder) LogFunc(logFunc func(sql string, args []any)) *MiddlewareBuilder {
	m.logFunc = logFunc
	return m
}
