package metrics

import (
	"context"
	"fmt"
	"sync"
	"time"

	"bamboo_lite_common_sdk/env"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
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

func GetOtelClient() *OtelClient {
	once.Do(func() {
		otelClient = &OtelClient{
			meter:      otel.GetMeterProvider().Meter(env.PSM()),
			counters:   make(map[string]metric.Int64Counter),
			histograms: make(map[string]metric.Float64Histogram),
		}
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
	hist := o.getOrCreateHistogram(name, "ms")
	hist.Record(ctx, latency.Seconds(), metric.WithAttributes(attrsFrom(tags...)...))
}

// attrsFrom Convert Tag slice to OTel attributes.
func attrsFrom(tags ...Tag) []attribute.KeyValue {
	if len(tags) == 0 {
		return nil
	}
	out := make([]attribute.KeyValue, 0, len(tags))
	for _, t := range tags {
		switch v := t.Value.(type) {
		case string:
			out = append(out, attribute.String(t.Name, v))
		case fmt.Stringer:
			out = append(out, attribute.String(t.Name, v.String()))
		case bool:
			out = append(out, attribute.Bool(t.Name, v))
		case int:
			out = append(out, attribute.Int(t.Name, v))
		case int32:
			out = append(out, attribute.Int64(t.Name, int64(v)))
		case int64:
			out = append(out, attribute.Int64(t.Name, v))
		case uint:
			out = append(out, attribute.Int64(t.Name, int64(v)))
		case uint32:
			out = append(out, attribute.Int64(t.Name, int64(v)))
		case uint64:
			out = append(out, attribute.Int64(t.Name, int64(v)))
		case float32:
			out = append(out, attribute.Float64(t.Name, float64(v)))
		case float64:
			out = append(out, attribute.Float64(t.Name, v))
		case []string:
			out = append(out, attribute.StringSlice(t.Name, v))
		default:
			out = append(out, attribute.String(t.Name, fmt.Sprint(v)))
		}
	}
	return out
}
