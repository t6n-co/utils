package metrics

import (
	"context"
	"time"
)

type ClientInterface interface {
	EmitCounter(ctx context.Context, name string, tags ...Tag)
	EmitLatency(ctx context.Context, name string, latency time.Duration, tags ...Tag)
}
