package metrics

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	requestCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "cart",
		Name:      "request_counter",
		Help:      "Counter of requests",
	}, []string{"url"})

	requestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "cart",
		Name:      "request_duration",
		Help:      "Duration of requests",
		Buckets:   prometheus.DefBuckets,
	}, []string{"url", "status"})

	externalRequestCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "cart",
		Name:      "external_request_counter",
		Help:      "Counter of external requests",
	}, []string{"url"})

	externalRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "cart",
		Name:      "external_request_duration",
		Help:      "Duration of external requests",
		Buckets:   prometheus.DefBuckets,
	}, []string{"url", "status"})

	repositoryAmounter = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "cart",
		Name:      "repository_amount",
		Help:      "Amount of objects in repository",
	}, []string{"repository"})
)

func RequestCounter(url string) {
	requestCounter.WithLabelValues(url).Inc()
}

func RequestDuration(url string, status int, duration float64) {
	requestDuration.WithLabelValues(url, strconv.Itoa(status)).Observe(duration)
}

func ExternalRequestCounter(url string) {
	externalRequestCounter.WithLabelValues(url).Inc()
}

func ExternalRequestDuration(url string, status string, duration float64) {
	externalRequestDuration.WithLabelValues(url, status).Observe(duration)
}

func ExternalRequestDurationWithError(url string, err error, duration float64) {
	status := "ok"
	if err != nil {
		status = "error"
	}
	ExternalRequestDuration(url, status, duration)
}

func CartRepositoryAmounter(amount float64) {
	repositoryAmounter.WithLabelValues("cart").Set(amount)
}
