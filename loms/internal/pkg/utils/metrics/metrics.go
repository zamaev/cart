package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	requestCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "loms",
		Name:      "request_counter",
		Help:      "Counter of requests",
	}, []string{"handler"})

	requestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "loms",
		Name:      "request_duration",
		Help:      "Duration of requests",
		Buckets:   prometheus.DefBuckets,
	}, []string{"handler", "status"})

	dbRequestCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "loms",
		Name:      "db_request_counter",
		Help:      "Counter of db requests",
	}, []string{"type"})

	dbRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "loms",
		Name:      "db_request_duration",
		Help:      "Duration of db requests",
		Buckets:   prometheus.DefBuckets,
	}, []string{"type", "status"})
)

func RequestCounter(handler string) {
	requestCounter.WithLabelValues(handler).Inc()
}

func RequestDuration(handler string, status string, duration float64) {
	requestDuration.WithLabelValues(handler, status).Observe(duration)
}

func RequestDurationWithErrorStart(handler string, err error, start time.Time) {
	status := "ok"
	if err != nil {
		status = "error"
	}
	RequestDuration(handler, status, time.Since(start).Seconds())
}

func DbRequestCounter(queryType string) {
	dbRequestCounter.WithLabelValues(queryType).Inc()
}

func DbRequestDuration(queryType string, status string, duration float64) {
	dbRequestDuration.WithLabelValues(queryType, status).Observe(duration)
}

func DbRequestDurationWithErrorStart(queryType string, err error, start time.Time) {
	status := "ok"
	if err != nil {
		status = "error"
	}
	DbRequestDuration(queryType, status, time.Since(start).Seconds())
}
