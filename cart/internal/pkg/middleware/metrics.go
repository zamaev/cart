package middleware

import (
	"net/http"
	"route256/cart/internal/pkg/utils/metrics"
	"time"
)

type MuxMetricsHandler interface {
	Handle(pattern string, handler http.Handler)
}

type MuxMetricsWrapper struct {
	MuxMetricsHandler
}

func NewMuxMetricsWrapper(mux MuxMetricsHandler) *MuxMetricsWrapper {
	return &MuxMetricsWrapper{MuxMetricsHandler: mux}
}

func (m *MuxMetricsWrapper) Handle(pattern string, handler http.Handler) {
	m.MuxMetricsHandler.Handle(pattern, HandlerMetricsWrapper{handler, pattern})
}

type HandlerMetricsWrapper struct {
	http.Handler
	urlPattern string
}

func (h HandlerMetricsWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ww := &responseWriterStatusWrapper{w, http.StatusOK}

	defer func(start time.Time) {
		metrics.RequestDuration(h.urlPattern, ww.status, time.Since(start).Seconds())
	}(time.Now())

	metrics.RequestCounter(h.urlPattern)

	h.Handler.ServeHTTP(ww, r)
}

type responseWriterStatusWrapper struct {
	http.ResponseWriter
	status int
}

func (w *responseWriterStatusWrapper) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
