package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func Set(trc trace.Tracer) {
	tracer = trc
}

func Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return tracer.Start(ctx, spanName, opts...)
}

func EndWithCheckError(span trace.Span, err *error) {
	if *err != nil {
		span.SetStatus(codes.Error, (*err).Error())
		span.RecordError(*err)
	}
	span.End()
}

func NewTracerProvider(ctx context.Context, serviceName, tracerUrl string) (*sdktrace.TracerProvider, error) {
	exporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL(tracerUrl))
	if err != nil {
		return nil, fmt.Errorf("otlptracehttp.New: %w", err)
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("resource.Merge: %w", err)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	return tracerProvider, nil
}
