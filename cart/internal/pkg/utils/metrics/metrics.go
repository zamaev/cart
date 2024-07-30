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

	cahceHitCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "cart",
		Name:      "cache_hit_counter",
		Help:      "Counter of cache hits",
	}, []string{"service_handler"})

	cacheHitDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "cart",
		Name:      "cache_hit_duration",
		Help:      "Duration of cache hits",
		Buckets:   prometheus.DefBuckets,
	}, []string{"service_handler"})

	cacheMissCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "cart",
		Name:      "cache_miss_counter",
		Help:      "Counter of cache misses",
	}, []string{"service_handler"})

	cacheMissDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "cart",
		Name:      "cache_miss_duration",
		Help:      "Duration of cache misses",
		Buckets:   prometheus.DefBuckets,
	}, []string{"service_handler"})
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

func CacheHitCounter(service_handler string) {
	cahceHitCounter.WithLabelValues(service_handler).Inc()
}

func CacheHitDuration(service_handler string, duration float64) {
	cacheHitDuration.WithLabelValues(service_handler).Observe(duration)
}

func CacheMissCounter(service_handler string) {
	cacheMissCounter.WithLabelValues(service_handler).Inc()
}

func CacheMissDuration(service_handler string, duration float64) {
	cacheMissDuration.WithLabelValues(service_handler).Observe(duration)
}
