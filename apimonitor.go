package apipulse

import (
	"net/http"
	"time"

	"github.com/kingztech2019/apipulse/pkg/health"
	"github.com/kingztech2019/apipulse/pkg/interceptor"
	"github.com/kingztech2019/apipulse/pkg/logging"
	"github.com/kingztech2019/apipulse/pkg/metrics"
	"github.com/kingztech2019/apipulse/pkg/rate"
	"github.com/kingztech2019/apipulse/pkg/recovery"
)

type APIMonitor struct {
	Metrics     *metrics.Metrics
	HealthCheck *health.HealthChecker
	RateLimiter *rate.RateLimiter
	Interceptor *interceptor.Interceptor
}

// The New function initializes and returns an APIMonitor instance with metrics, health checker, rate
// limiter, and interceptors configured.
func New(limit int, interval time.Duration) *APIMonitor {
	m := metrics.NewMetrics()
	hc := health.NewHealthChecker()
	rl := rate.NewRateLimiter(limit, interval)

	i := interceptor.NewInterceptor(
		logging.LoggingMiddleware,
		m.MetricsMiddleware,
		recovery.ErrorTrackingMiddleware,
		rl.RateLimitMiddleware,
		// metrics.LatencyMiddleware,
	)

	return &APIMonitor{
		Metrics:     m,
		HealthCheck: hc,
		RateLimiter: rl,
		Interceptor: i,
	}
}

func (m *APIMonitor) RegisterMetrics() {
	m.Metrics.RegisterMetrics()
}

func (m *APIMonitor) MetricsHandler() http.Handler {
	return m.Metrics.MetricsHandler()
}

func (m *APIMonitor) HealthCheckHandler() http.Handler {
	return m.HealthCheck.HealthCheckHandler()
}
