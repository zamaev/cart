package middleware

import (
	"errors"
	"log/slog"
	"net/http"
	"route256/cart/internal/pkg/customerror"
)

type ErrorWrapper func(w http.ResponseWriter, r *http.Request) error

func (h ErrorWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		slog.Error("Handle error", "method", r.Method, "uri", r.URL.Path, "err", err)

		var errStatusCode customerror.ErrStatusCode
		if errors.As(err, &errStatusCode) {
			w.WriteHeader(errStatusCode.Status)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Write([]byte("{}"))
	}
}
