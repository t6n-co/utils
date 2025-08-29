package logCtx

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/t6n-co/utils/env"
	"github.com/t6n-co/utils/metrics"
)

func formatString(fmtString string, args ...interface{}) string {
	return fmt.Sprintf(fmtString, args...)
}

func getBaseArgs(ctx context.Context) []any {
	return []any{"psm", env.PSM()}
}

func Info(ctx context.Context, fmtString string, args ...interface{}) {
	slog.Info(formatString(fmtString, args...), getBaseArgs(ctx)...)
}

func Warn(ctx context.Context, fmtString string, args ...interface{}) {
	metrics.GetClient().EmitCounter(ctx, "log.warn")
	slog.Warn(formatString(fmtString, args...), getBaseArgs(ctx)...)
}

func Error(ctx context.Context, fmtString string, args ...interface{}) {
	metrics.GetClient().EmitCounter(ctx, "log.error")
	slog.Error(formatString(fmtString, args...), getBaseArgs(ctx)...)
}

func Debug(ctx context.Context, fmtString string, args ...interface{}) {
	slog.Debug(formatString(fmtString, args...), getBaseArgs(ctx)...)
}
