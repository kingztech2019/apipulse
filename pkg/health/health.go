package health

import (
	"encoding/json"
	"net/http"
	"sync"
)

type CustomHealthCheck func() bool

type HealthStatus struct {
	Status  string            `json:"status"`
	Details map[string]string `json:"details"`
}

type HealthChecker struct {
	checks map[string]CustomHealthCheck
	mu     sync.Mutex
}

func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		checks: make(map[string]CustomHealthCheck),
	}
}

// The `RegisterCheck` method in the `HealthChecker` struct is used to register a custom health check
// function with a given name.
func (hc *HealthChecker) RegisterCheck(name string, check CustomHealthCheck) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	hc.checks[name] = check
}

// This function `HealthCheckHandler` is defining an HTTP handler function that will handle incoming
// HTTP requests related to health checks. Here's a breakdown of what it does:
func (hc *HealthChecker) HealthCheckHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		status := HealthStatus{
			Status:  "healthy",
			Details: make(map[string]string),
		}
		hc.mu.Lock()
		defer hc.mu.Unlock()
		for name, check := range hc.checks {
			if !check() {
				status.Status = "unhealthy"
				status.Details[name] = "unhealthy"
			} else {
				status.Details[name] = "healthy"
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
	})
}
