package opentelemetry

import (
	"context"
	"fmt"
	gsql "github.com/DaHuangQwQ/gsql"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const instrumentationName = "gitee.com/geektime-geekbang/geektime-go/gsql/middlewares/opentelemetry"

type MiddlewareBuilder struct {
	Tracer trace.Tracer
}

func (m MiddlewareBuilder) Build() gsql.Middleware {
	if m.Tracer == nil {
		m.Tracer = otel.GetTracerProvider().Tracer(instrumentationName)
	}
	return func(next gsql.Handler) gsql.Handler {
		return func(ctx context.Context, qc *gsql.QueryContext) *gsql.QueryResult {
			tbl := qc.Model.TableName
			spanCtx, span := m.Tracer.Start(ctx, fmt.Sprintf("%s-%s", qc.Type, tbl))
			defer span.End()

			q, _ := qc.Builder.Build()
			if q != nil {
				span.SetAttributes(attribute.String("sql", q.SQL))
			}

			span.SetAttributes(attribute.String("table", tbl))
			span.SetAttributes(attribute.String("component", "gsql"))
			res := next(spanCtx, qc)
			if res.Err != nil {
				span.RecordError(res.Err)
			}
			return res
		}
	}
}
