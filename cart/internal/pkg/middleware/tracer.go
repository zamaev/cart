package middleware

import (
	"net/http"
	"route256/cart/pkg/logger"
	"route256/cart/pkg/tracing"
)

type MuxTracerHandler interface {
	Handle(pattern string, handler http.Handler)
}

type MuxTracerWrapper struct {
	MuxTracerHandler
}

func NewMuxTracerWrapper(mux MuxTracerHandler) *MuxTracerWrapper {
	return &MuxTracerWrapper{mux}
}

func (m *MuxTracerWrapper) Handle(pattern string, handler http.Handler) {
	m.MuxTracerHandler.Handle(pattern, HandlerTracerWrapper{handler, pattern})
}

type HandlerTracerWrapper struct {
	handler    http.Handler
	urlPattern string
}

func (h HandlerTracerWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.Start(r.Context(), h.urlPattern)
	defer span.End()

	l := logger.With("trace_id", span.SpanContext().TraceID())
	ctx = logger.ToContext(ctx, l)

	r = r.WithContext(ctx)
	h.handler.ServeHTTP(w, r)
}
