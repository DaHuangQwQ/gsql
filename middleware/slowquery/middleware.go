package slowquery

import (
	"context"
	"github.com/DaHuangQwQ/gsql"
	"log"
	"time"
)

type MiddlewareBuilder struct {
	threshold time.Duration
	logFunc   func(query string, args []any)
}

func NewMiddlewareBuilder(threshold time.Duration) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		logFunc: func(query string, args []any) {
			log.Printf("sql: %s, args: %v", query, args)
		},
		threshold: threshold,
	}
}

func (m *MiddlewareBuilder) LogFunc(fn func(query string, args []any)) *MiddlewareBuilder {
	m.logFunc = fn
	return m
}

func (m MiddlewareBuilder) Build() gsql.Middleware {
	return func(next gsql.Handler) gsql.Handler {
		return func(ctx context.Context, qc *gsql.QueryContext) *gsql.QueryResult {
			startTime := time.Now()
			defer func() {
				duration := time.Since(startTime)
				if duration <= m.threshold {
					return
				}
				q, err := qc.Builder.Build()
				if err == nil {
					m.logFunc(q.SQL, q.Args)
				}
			}()

			return next(ctx, qc)
		}
	}
}
