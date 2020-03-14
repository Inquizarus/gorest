package gorest

import (
	"fmt"
	"net/http"
)

// Middleware for HTTP requests
type Middleware func(http.HandlerFunc) http.HandlerFunc

// WithJSONContent is a Middleware that sets the content type to json for the response
func WithJSONContent() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("content-type", "application/json")
			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}

// WithStrictTransportSecurity injects strict TLS related headers
func WithStrictTransportSecurity(TLS ServeTLSConfig) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
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
			f(w, r)
		}
	}
}

// ChainMiddleware creates a function that will apply all passed middlewares
func ChainMiddleware(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}
