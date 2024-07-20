package middleware

import (
	"context"
	"route256/loms/pkg/logger"
	"route256/loms/pkg/tracing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Panic(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	ctx, span := tracing.Start(ctx, "middleware.Panic")
	defer tracing.EndWithCheckError(span, &err)

	defer func() {
		if e := recover(); e != nil {
			span.AddEvent("panic")
			logger.Errorw(ctx, "panic", "err", e)
			err = status.Errorf(codes.Internal, "panic: %v", e)
		}
	}()
	resp, err = handler(ctx, req)
	return resp, err
}
