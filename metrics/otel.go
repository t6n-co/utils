package metrics

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Panda-Education/bamboo_lite_common_sdk/env"
	"github.com/Panda-Education/bamboo_lite_common_sdk/logCtx"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	metric2 "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

var (
	once       sync.Once
	otelClient *OtelClient
)

type OtelClient struct {
	mu         sync.RWMutex
	meter      metric.Meter
	counters   map[string]metric.Int64Counter
	histograms map[string]metric.Float64Histogram
}

func (o *OtelClient) getMeterProvider() *metric2.MeterProvider {
	exporter, err := prometheus.New(
		prometheus.WithNamespace(strings.NewReplacer(".", "_").Replace(env.PSM())),
	)
	if err != nil {
		panic(err)
	}
	meterProvider := metric2.NewMeterProvider(
		metric2.WithResource(resource.Default()),
		metric2.WithReader(exporter),
	)
	return meterProvider
}

func GetOtelClient() *OtelClient {
	once.Do(func() {
		otelClient = &OtelClient{
			counters:   make(map[string]metric.Int64Counter),
			histograms: make(map[string]metric.Float64Histogram),
		}
		meterProvider := otelClient.getMeterProvider()
		otelClient.meter = meterProvider.Meter(env.PSM())
	})
	return otelClient
}

func (o *OtelClient) getOrCreateCounter(name string) metric.Int64Counter {
	o.mu.RLock()
	if ctr, ok := o.counters[name]; ok {
		o.mu.RUnlock()
		return ctr
	}

	o.mu.RUnlock()
	o.mu.Lock()
	defer o.mu.Unlock()

	if c, ok := o.counters[name]; ok { // re-check under write lock
		return c
	}
	c, err := o.meter.Int64Counter(
		name,
		metric.WithUnit("1"),
	)
	if err != nil {
		panic(fmt.Errorf("create counter %q: %w", name, err))
	}
	o.counters[name] = c
	return c
}

func (o *OtelClient) getOrCreateHistogram(name, unit string) metric.Float64Histogram {
	o.mu.RLock()
	if h, ok := o.histograms[name]; ok {
		o.mu.RUnlock()
		return h
	}

	o.mu.RUnlock()
	o.mu.Lock()
	defer o.mu.Unlock()
	if h, ok := o.histograms[name]; ok {
		return h
	}

	h, err := o.meter.Float64Histogram(
		name,
		metric.WithUnit(unit),
	)
	if err != nil {
		panic(fmt.Errorf("create histogram %q: %w", name, err))
	}
	o.histograms[name] = h
	return h
}

func (o *OtelClient) EmitCounter(ctx context.Context, name string, tags ...Tag) {
	ctr := o.getOrCreateCounter(name)
	ctr.Add(ctx, 1, metric.WithAttributes(attrsFrom(tags...)...))
}

func (o *OtelClient) EmitLatency(ctx context.Context, name string, latency time.Duration, tags ...Tag) {
	hist := o.getOrCreateHistogram(name, "s")
	hist.Record(ctx, latency.Seconds(), metric.WithAttributes(attrsFrom(tags...)...))
}

func (o *OtelClient) ExposeToPrometheus(ctx context.Context, opt ExposeToPrometheusOpt) {
	if opt.Endpoint == "" {
		opt.Endpoint = "/metrics"
	}
	if opt.Host == "" {
		opt.Host = "0.0.0.0"
	}
	if opt.Port == 0 {
		opt.Port = 8080
	}

	mux := http.NewServeMux()
	mux.Handle(opt.Endpoint, promhttp.Handler())

	addr := fmt.Sprintf("%s:%d", opt.Host, opt.Port)
	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	logCtx.Info(ctx, "[OtelClient.ExposeToPrometheus] Exposing metrics on http://%s%s", addr, opt.Endpoint)

	// Start server
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logCtx.Error(ctx, "[OtelClient.ExposeToPrometheus] Prometheus server error: %v", err)
		}
	}()

	// Graceful shutdown on ctx cancel
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			logCtx.Error(ctx, "[OtelClient.ExposeToPrometheus] Prometheus server shutdown error: %v", err)
		}
	}()
}

func (o *OtelClient) GetPrometheusHandler(ctx context.Context) http.Handler {
	return promhttp.Handler()
}
