package middleware

/*
Goals:
- log the total number of requests that have been served
- log the level of concurrency at which request are served (number of concurrent requests)
- log the latency characteristics of total request lifetime
- log the latency characteristics of individual operations
	- read mapping from database
	- write mapping to database
	- read mapping from cache
	- write mapping to cache
	- generate random 62 bit id
-
*/

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	TotalRequests prometheus.Counter
	ConcurrentRequests prometheus.Gauge
	RequestLatency prometheus.Histogram
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	metrics := &Metrics{
		TotalRequests: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "total_requests",
				Help: "total number of requests served across all endpoints by this server",
			},
		),
		ConcurrentRequests: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "concurrent_requests",
				Help: "gauge of concurrently running requests",
			},
		),
		RequestLatency: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Name: "request_latency",
				Help: "the histogram of request latencies for all routes",
			},
		),
	}

	reg.MustRegister(metrics.TotalRequests)
	reg.MustRegister(metrics.ConcurrentRequests)
	reg.MustRegister(metrics.RequestLatency)
	return metrics
}

func MetricsMiddleware(metrics *Metrics, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metrics.TotalRequests.Add(1.0)
		metrics.ConcurrentRequests.Inc()
		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start).Seconds()
		metrics.ConcurrentRequests.Dec()
		metrics.RequestLatency.Observe(duration)
	})
}