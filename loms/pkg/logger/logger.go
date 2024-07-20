package logger

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerCtxKey struct{}

type Logger struct {
	logger *zap.SugaredLogger
}

var globalLogger *Logger

func init() {
	var err error
	globalLogger, err = New()
	if err != nil {
		panic(err)
	}
}

func Set(logger *Logger) {
	globalLogger = logger
}

func New() (*Logger, error) {
	config := zap.NewProductionConfig()
	config.Level.SetLevel(zapcore.InfoLevel)
	config.DisableCaller = true

	logger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("config.Build: %w", err)
	}

	return &Logger{logger: logger.Sugar()}, nil
}

func With(args ...any) *Logger {
	return &Logger{globalLogger.logger.With(args...)}
}

func ToContext(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey{}, logger)
}

func Panicw(ctx context.Context, msg string, keyAndValues ...any) {
	spanID := trace.SpanFromContext(ctx).SpanContext().SpanID()

	if loggerCtx, ok := ctx.Value(loggerCtxKey{}).(*Logger); ok {
		loggerCtx.logger.With("spanID", spanID).Panicw(msg, keyAndValues...)
		return
	}
	globalLogger.logger.With("spanID", spanID).Panicw(msg, keyAndValues...)
}

func Errorw(ctx context.Context, msg string, keyAndValues ...any) {
	spanID := trace.SpanFromContext(ctx).SpanContext().SpanID()

	if loggerCtx, ok := ctx.Value(loggerCtxKey{}).(*Logger); ok {
		loggerCtx.logger.With("spanID", spanID).Errorw(msg, keyAndValues...)
		return
	}
	globalLogger.logger.With("spanID", spanID).Errorw(msg, keyAndValues...)
}

func Infow(ctx context.Context, msg string, keyAndValues ...any) {
	spanID := trace.SpanFromContext(ctx).SpanContext().SpanID()

	if loggerCtx, ok := ctx.Value(loggerCtxKey{}).(*Logger); ok {
		loggerCtx.logger.With("spanID", spanID).Infow(msg, keyAndValues...)
		return
	}
	globalLogger.logger.With("spanID", spanID).Infow(msg, keyAndValues...)
}
