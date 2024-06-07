package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	RequestCount    *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec
	RequestErrors   *prometheus.CounterVec
	InProgress      prometheus.Gauge
	RequestSize     *prometheus.HistogramVec
	ResponseSize    *prometheus.HistogramVec
}

type responseWriter struct {
	http.ResponseWriter
	status      int
	writtenSize int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(data)
	rw.writtenSize += size
	return size, err
}

func WrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK, 0}
}

// The NewMetrics function initializes and returns a struct containing various Prometheus metrics for
// monitoring API requests.
func NewMetrics() *Metrics {
	return &Metrics{
		RequestCount: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "api_request_count",
				Help: "Total number of API requests",
			},
			[]string{"path", "method", "status_code"},
		),
		RequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "api_request_duration_seconds",
				Help:    "Histogram of API request durations",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"path", "method"},
		),
		RequestErrors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "api_request_errors",
				Help: "Total number of API request errors",
			},
			[]string{"path", "method"},
		),
		InProgress: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "api_in_progress_requests",
				Help: "Current number of in-progress requests",
			},
		),
		RequestSize: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "api_request_size_bytes",
				Help:    "Histogram of API request sizes",
				Buckets: prometheus.ExponentialBuckets(100, 10, 5),
			},
			[]string{"path", "method"},
		),
		ResponseSize: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "api_response_size_bytes",
				Help:    "Histogram of API response sizes",
				Buckets: prometheus.ExponentialBuckets(100, 10, 5),
			},
			[]string{"path", "method"},
		),
	}
}

func (m *Metrics) RegisterMetrics() {
	prometheus.MustRegister(m.RequestCount)
	prometheus.MustRegister(m.RequestDuration)
	prometheus.MustRegister(m.RequestErrors)
	prometheus.MustRegister(m.InProgress)
	prometheus.MustRegister(m.RequestSize)
	prometheus.MustRegister(m.ResponseSize)
}

// The `func (m *Metrics) MetricsMiddleware(next http.Handler) http.Handler` function is a middleware
// function that wraps an HTTP handler function. It captures metrics related to incoming HTTP requests
// and responses.
func (m *Metrics) MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := WrapResponseWriter(w)
		m.InProgress.Inc()
		defer m.InProgress.Dec()

		reqSize := float64(r.ContentLength)
		if reqSize < 0 {
			reqSize = 0
		}

		next.ServeHTTP(rw, r)
		duration := time.Since(start).Seconds()

		m.RequestCount.WithLabelValues(r.URL.Path, r.Method, http.StatusText(rw.status)).Inc()
		m.RequestDuration.WithLabelValues(r.URL.Path, r.Method).Observe(duration)
		m.RequestSize.WithLabelValues(r.URL.Path, r.Method).Observe(reqSize)
		m.ResponseSize.WithLabelValues(r.URL.Path, r.Method).Observe(float64(rw.writtenSize))

		if rw.status >= 400 {
			m.RequestErrors.WithLabelValues(r.URL.Path, r.Method).Inc()
		}
	})
}

func (m *Metrics) MetricsHandler() http.Handler {
	return promhttp.Handler()
}
