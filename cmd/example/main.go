package main

import (
	"net/http"
	"time"

	"github.com/kingztech2019/apipulse"
)

func main() {
	monitor := apipulse.New(1, time.Minute)
	monitor.RegisterMetrics()

	mux := http.NewServeMux()
	mux.Handle("/metrics", monitor.MetricsHandler())

	mux.Handle("/health", monitor.HealthCheckHandler())
	// Use Interceptor to apply middleware stack to all routes
	mux.Handle("/", monitor.Interceptor.Intercept(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte("Hello, World!"))
	})))

	// Add additional routes
	mux.Handle("/example", monitor.Interceptor.Intercept(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte("Example Route"))
	})))

	http.ListenAndServe(":8080", mux)
}
