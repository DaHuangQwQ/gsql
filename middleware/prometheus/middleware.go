package prometheus

import (
	"context"
	"github.com/DaHuangQwQ/gsql"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type MiddlewareBuilder struct {
	Namespace   string
	Subsystem   string
	Name        string
	Help        string
	ConstLabels map[string]string
}

func (m MiddlewareBuilder) Build() gsql.Middleware {
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:        m.Name,
		Subsystem:   m.Subsystem,
		Namespace:   m.Name,
		Help:        m.Help,
		ConstLabels: m.ConstLabels,
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.90:  0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	}, []string{"type", "table"})

	prometheus.MustRegister(vector)

	return func(next gsql.Handler) gsql.Handler {
		return func(ctx context.Context, qc *gsql.QueryContext) *gsql.QueryResult {
			startTime := time.Now()
			defer func() {
				// 执行时间
				vector.WithLabelValues(qc.Type.String(), qc.Model.TableName).Observe(float64(time.Since(startTime).Milliseconds()))
			}()
			return next(ctx, qc)
		}
	}
}
