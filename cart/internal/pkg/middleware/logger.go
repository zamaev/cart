package middleware

import (
	"net/http"
	"route256/cart/pkg/logger"
	"time"
)

// Deprecated: use middleware.MetricsWrapper
type LoggerWrapperHandler struct {
	Wrap http.Handler
}

func (h LoggerWrapperHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	logger.Infow(r.Context(), "request started", "method", r.Method, "url", r.URL.Path)
	h.Wrap.ServeHTTP(w, r)
	logger.Infow(r.Context(), "request ended", "method", r.Method, "url", r.URL.Path, "duration", time.Since(start))
}
