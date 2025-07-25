package metrics

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "proxy_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path"},
	)
	ErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "proxy_errors_total",
			Help: "Total number of HTTP errors",
		},
		[]string{"method", "path", "code"},
	)
	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "proxy_request_duration_seconds",
			Help:    "Request latency",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

func Init() {
	prometheus.MustRegister(RequestsTotal, ErrorsTotal, RequestDuration)
}

func Handler() http.Handler {
	return promhttp.Handler()
}

func ObserveRequest(method, path string, code int, duration time.Duration) {
	RequestsTotal.WithLabelValues(method, path).Inc()
	RequestDuration.WithLabelValues(method, path).Observe(duration.Seconds())
	if code >= 400 {
		ErrorsTotal.WithLabelValues(method, path, itoa(code)).Inc()
	}
}

func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}
