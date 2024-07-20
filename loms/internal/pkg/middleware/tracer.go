package middleware

import (
	"context"
	"route256/loms/pkg/logger"
	"route256/loms/pkg/tracing"

	"google.golang.org/grpc"
)

func Tracer(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ any, err error) {
	ctx, span := tracing.Start(ctx, info.FullMethod)
	defer tracing.EndWithCheckError(span, &err)

	l := logger.With("trace_id", span.SpanContext().TraceID())
	ctx = logger.ToContext(ctx, l)

	return handler(ctx, req)
}
