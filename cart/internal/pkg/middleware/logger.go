package middleware

import (
	"log"
	"net/http"
	"time"
)

type LoggerWrapperHandler struct {
	Wrap http.Handler
}

func (h LoggerWrapperHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	log.Printf("%s %s start\n", r.Method, r.URL.Path)
	h.Wrap.ServeHTTP(w, r)
	log.Printf("%s %s ended. Duration: %s\n", r.Method, r.URL.Path, time.Since(start))
}
