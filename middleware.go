package gorest

import (
	"fmt"
	"net/http"
)

// Middleware for HTTP requests
type Middleware func(http.Handler) http.Handler

// WithJSONContent is a Middleware that sets the content type to json for the response
func WithJSONContent() Middleware {
	return func(f http.Handler) http.Handler {
		// Define the http.HandlerFunc
		fn := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("content-type", "application/json")
			// Call the next middleware/handler in chain
			f.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// WithStrictTransportSecurity injects strict TLS related headers
func WithStrictTransportSecurity(TLS ServeTLSConfig) Middleware {
	return func(f http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if true == TLS.Strict {
				headerValue := "max-age=63072000"
				if true == TLS.StrictSubDomains {
					headerValue = fmt.Sprintf("%s; includeSubDomains", headerValue)
				}
				if true == TLS.Preload {
					headerValue = fmt.Sprintf("%s; preload", headerValue)
				}
				w.Header().Add("Strict-Transport-Security", headerValue)
			}
			f.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// ChainMiddleware creates a function that will apply all passed middlewares
func ChainMiddleware(f http.Handler, middlewares ...Middleware) http.Handler {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}
