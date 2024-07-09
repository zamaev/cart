package middleware

import (
	"errors"
	"net/http"
	"route256/cart/internal/pkg/customerror"
	"route256/cart/pkg/logger"
)

type ErrorWrapper func(w http.ResponseWriter, r *http.Request) error

func (h ErrorWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		logger.Errorw(r.Context(), "Handle error", "method", r.Method, "url", r.URL.Path, "err", err)

		var errStatusCode customerror.ErrStatusCode
		if errors.As(err, &errStatusCode) {
			w.WriteHeader(errStatusCode.Status)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Write([]byte("{}"))
	}
}
