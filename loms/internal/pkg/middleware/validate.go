package middleware

import (
	"context"
	"route256/loms/pkg/tracing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Validate(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ any, err error) {
	ctx, span := tracing.Start(ctx, "middleware.Validate")
	defer tracing.EndWithCheckError(span, &err)

	if v, ok := req.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}
	return handler(ctx, req)
}
