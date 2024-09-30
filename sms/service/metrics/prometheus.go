package metrics

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"time"
	"webook/sms/service"
)

type Decorator struct {
	svc    service.Service
	vector *prometheus.SummaryVec
}

func NewDecorator(svc service.Service) *Decorator {
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "DaHuang",
		Subsystem: "webook",
		Name:      "sms_resp_time",
		Help:      "统计 SMS 服务性能数据",
	}, []string{"tpl_id"})
	prometheus.MustRegister(vector)
	return &Decorator{
		svc:    svc,
		vector: vector,
	}
}

func (d *Decorator) Send(ctx context.Context,
	tplId string, args []service.NamedArg, numbers ...string) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		d.vector.WithLabelValues(tplId).Observe(float64(duration))
	}()
	return d.svc.Send(ctx, tplId, args, numbers...)
}
