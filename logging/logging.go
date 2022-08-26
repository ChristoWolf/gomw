// Package logging implements a simple logger middleware.
package logging

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

// wrappedResponseWriter is a wrapper around
// http.ResponseWriter that caches additional
// logging relevant information.
type wrappedResponseWriter struct {
	http.ResponseWriter
	body   []byte
	status int
}

// NewWrappedResponseWriter creates a new
// wrappedResponseWriter with overwritable
// default values which will be filled with
// request execution information.
func NewWrappedResponseWriter(w http.ResponseWriter) *wrappedResponseWriter {
	return &wrappedResponseWriter{
		ResponseWriter: w,
		body:           []byte{},
		status:         http.StatusInternalServerError}
}

// WriteHeader partially implements http.ResponseWriter.
// It caches the status code and calls the underlying
// http.ResponseWriter.WriteHeader.
func (w *wrappedResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// Write partially implements http.ResponseWriter.
// It caches the response body and calls the underlying
// http.ResponseWriter.Write.
func (w *wrappedResponseWriter) Write(b []byte) (int, error) {
	w.body = b
	return w.ResponseWriter.Write(b)
}

// Middleware is a middleware for logging information
// about the request and response.
// TODO: Log more information (sizes, execution time, etc).
func Middleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			wrapped := NewWrappedResponseWriter(w)
			body, err := io.ReadAll(r.Body)
			logger.Printf("Request: %s %s %s\n", r.Method, r.URL, body)
			if err != nil {
				err := fmt.Errorf("error reading request body: %w", err)
				logger.Printf("%v\n", err)
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprint(w, err)
				}).ServeHTTP(wrapped, r)
			} else {
				next.ServeHTTP(wrapped, r)
			}
			logger.Printf("Response: %d %s\n", wrapped.status, string(wrapped.body))
		})
	}
}
