 # APIPulse

APIPulse is a comprehensive Go package designed to monitor, log, and manage API requests. It provides essential middleware functionalities including metrics collection, request logging, health checks, rate limiting, and error recovery, making it easier to ensure the reliability, performance, and security of your API.


## Installation

To install the package, run:

```sh
go get github.com/kingzmentech2019/apipulse
```
## Features
* **Metrics Collection:** Gather detailed metrics on API request counts, durations, errors, sizes, and more using Prometheus.

* **Request Logging:** Log incoming requests and responses with their statuses.

* **Health Checks:** Define custom health checks and provide a health status endpoint.

* **Rate Limiting:** Limit the rate of incoming requests based on client IP addresses.

* **Error Recovery:** Recover from panics in HTTP handlers and log the errors.

# Usage

## Basic Setup
Here’s an example of how to set up a basic API server using the APIMonitor package.

```go
package main

import (
    "net/http"
    "time"

    "github.com/kingzmentech2019/apipulse"
)

func main() {
    monitor := apipulse.New(100, time.Minute)
    monitor.RegisterMetrics()

    mux := http.NewServeMux()
    mux.Handle("/metrics", monitor.MetricsHandler())
    mux.Handle("/health", monitor.HealthCheckHandler())

    // Add middleware to the root handler
    mux.Handle("/", monitor.Interceptor.Intercept(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })))

    http.ListenAndServe(":8080", mux)
}

```
## Features
1. Metrics Collection

APIPulse collects and exposes metrics in a format compatible with Prometheus.

* **RegisterMetrics:** Registers the necessary metrics with Prometheus.
* **MetricsHandler:** An HTTP handler that exposes the metrics endpoint.
#### Example

```go
monitor := apipulse.New(100, time.Minute)
monitor.RegisterMetrics()

mux.Handle("/metrics", monitor.MetricsHandler())

```
### 2. Health Checks

APIPulse allows you to define health checks for various parts of your application.

* **HealthCheckHandler:** An HTTP handler that exposes the health check endpoint.
* **RegisterCheck:** Registers a new health check.
#### Example

```go
monitor := apipulse.New(100, time.Minute)
monitor.HealthCheck.RegisterCheck("database", func() bool {
    // Add logic to check database health
    return true
})

mux.Handle("/health", monitor.HealthCheckHandler())
```
### 3. Request Logging

APIPulse logs each request made to your API.

* **LoggingMiddleware:** Middleware for logging requests.
#### Example

```go
mux.Handle("/", monitor.Interceptor.Intercept(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello, World!"))
})))

```
### 4. Error Tracking

APIPulse tracks and logs errors occurring in your API.

* **ErrorTrackingMiddleware:** Middleware for tracking errors.
#### Example

```go
mux.Handle("/", monitor.Interceptor.Intercept(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    panic("something went wrong")
})))

```
### 5. Rate Limiting

APIPulse allows you to limit the rate of incoming requests.

* **RateLimitMiddleware:** Middleware for rate limiting.

#### Example

```go
monitor := apipulse.New(100, time.Minute)

mux.Handle("/", monitor.Interceptor.Intercept(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello, World!"))
})))

```
### Full Example

Here is a complete example that demonstrates how to use the APIPulse package:

```go
package main

import (
    "net/http"
    "time"

    "github.com/kingzmentech2019/apipulse"
)

func main() {
    monitor := apipulse.New(100, time.Minute)
    monitor.RegisterMetrics()

    mux := http.NewServeMux()
    mux.Handle("/metrics", monitor.MetricsHandler())
    mux.Handle("/health", monitor.HealthCheckHandler())

    // Add additional routes and apply middleware
    mux.Handle("/", monitor.Interceptor.Intercept(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })))

    mux.Handle("/example", monitor.Interceptor.Intercept(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Example Route"))
    })))

    http.ListenAndServe(":8080", mux)
}

```