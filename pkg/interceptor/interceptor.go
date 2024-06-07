package interceptor

import (
	"net/http"
)

// Middleware defines the function signature for middleware
type Middleware func(http.Handler) http.Handler

// Interceptor allows chaining middleware
type Interceptor struct {
	middlewares []Middleware
}

// NewInterceptor creates a new Interceptor instance
func NewInterceptor(middlewares ...Middleware) *Interceptor {
	return &Interceptor{middlewares: middlewares}
}

// Use adds a new middleware to the stack
func (i *Interceptor) Use(middleware Middleware) {
	i.middlewares = append(i.middlewares, middleware)
}

// The `Intercept` method in the `Interceptor` struct is iterating over the list of middlewares in
// reverse order and applying each middleware to the provided `http.Handler`. It starts from the last
// middleware added and works its way to the first one added. Each middleware is applied to the
// handler, modifying or augmenting its behavior. Finally, the modified handler is returned after all
// the middlewares have been applied. This allows for chaining multiple middlewares together to
// intercept and process HTTP requests and responses.
func (i *Interceptor) Intercept(handler http.Handler) http.Handler {
	for j := len(i.middlewares) - 1; j >= 0; j-- {
		handler = i.middlewares[j](handler)
	}
	return handler
}
