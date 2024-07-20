package middleware

import (
	"context"
	"route256/loms/internal/pkg/utils/metrics"
	"route256/loms/pkg/tracing"
	"time"

	"google.golang.org/grpc"
)

func Metrics(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ any, err error) {
	ctx, span := tracing.Start(ctx, "middleware.Metrics")
	defer tracing.EndWithCheckError(span, &err)

	metrics.RequestCounter(info.FullMethod)
	defer func(start time.Time) {
		metrics.RequestDurationWithErrorStart(info.FullMethod, err, start)
	}(time.Now())

	return handler(ctx, req)
}
