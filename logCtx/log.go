package logCtx

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Panda-Education/bamboo_lite_common_sdk/env"
)

func formatString(fmtString string, args ...interface{}) string {
	return fmt.Sprintf(fmtString, args...)
}

func getBaseArgs(ctx context.Context) []interface{} {
	var args []interface{}
	args = append(args, []string{"PSM", env.PSM()})
	return args
}

func Info(ctx context.Context, fmtString string, args ...interface{}) {
	slog.Info(formatString(fmtString, args...), getBaseArgs(ctx)...)
}

func Warn(ctx context.Context, fmtString string, args ...interface{}) {
	slog.Warn(formatString(fmtString, args...), getBaseArgs(ctx)...)
}

func Error(ctx context.Context, fmtString string, args ...interface{}) {
	slog.Error(formatString(fmtString, args...), getBaseArgs(ctx)...)
}

func Debug(ctx context.Context, fmtString string, args ...interface{}) {
	slog.Debug(formatString(fmtString, args...), getBaseArgs(ctx)...)
}
