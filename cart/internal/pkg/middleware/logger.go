package middleware

import (
	"log"
	"net/http"
)

type LoggerWrapperHandler struct {
	Wrap http.Handler
}

func (h LoggerWrapperHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s\n", r.Method, r.URL.Path)
	h.Wrap.ServeHTTP(w, r)
}
