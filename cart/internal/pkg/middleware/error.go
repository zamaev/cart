package middleware

import (
	"log/slog"
	"net/http"
)

type ErrorWrapper func(w http.ResponseWriter, r *http.Request) error

func (h ErrorWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		slog.Error("Handle error", "method", r.Method, "uri", r.URL.Path, "err", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{}"))
	}
}
