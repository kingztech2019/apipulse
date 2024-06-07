package logging

import (
	"log"
	"net/http"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func WrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

// The `LoggingMiddleware` function logs incoming requests and processed responses in Go.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Incoming request: Method=%s, URL=%s, Headers=%v", r.Method, r.URL, r.Header)
		rw := WrapResponseWriter(w)
		next.ServeHTTP(rw, r)
		log.Printf("Request processed: Status=%d", rw.status)
	})
}
