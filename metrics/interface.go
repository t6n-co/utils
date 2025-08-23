package metrics

import (
	"context"
	"net/http"
	"time"
)

type ExposeToPrometheusOpt struct {
	Port     int
	Host     string
	Endpoint string
}

type ClientInterface interface {
	EmitCounter(ctx context.Context, name string, tags ...Tag)
	EmitLatency(ctx context.Context, name string, latency time.Duration, tags ...Tag)
	ExposeToPrometheus(ctx context.Context, opt ExposeToPrometheusOpt)
	GetPrometheusHandler(ctx context.Context) http.Handler
}
