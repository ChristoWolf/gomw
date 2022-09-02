// Package logging implements a simple logger middleware.
package logging

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
	"time"
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
func Middleware(options Options) func(http.Handler) http.Handler {
	reqTemplate := template.Must(template.New("request").Parse(requestTemplate))
	respTemplate := template.Must(template.New("response").Parse(responseTemplate))
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			wrapped := NewWrappedResponseWriter(w)
			var reqBody []byte
			var err error
			if options.withBodies {
				reqBody, err = io.ReadAll(r.Body)
				if err != nil {
					err := fmt.Errorf("error reading request body: %w", err)
					options.logger.Printf("Error: %v\n", err)
					http.Error(wrapped, err.Error(), http.StatusInternalServerError)
					return
				}
			}
			reqData := requestData{
				WithBody:          options.withBodies,
				WithMethod:        options.withMethod,
				WithUrl:           options.withUrl,
				WithContentLength: options.withContentLengths,
				Body:              string(reqBody),
				Method:            r.Method,
				Url:               r.URL.String(),
				ContentLength:     r.ContentLength,
			}
			writer := strings.Builder{}
			err = reqTemplate.Execute(&writer, reqData)
			if err != nil {
				err := fmt.Errorf("error executing request template: %w", err)
				options.logger.Printf("Error: %v\n", err)
				http.Error(wrapped, err.Error(), http.StatusInternalServerError)
				return
			}
			options.logger.Println(writer.String())
			start := time.Now()
			next.ServeHTTP(wrapped, r)
			duration := time.Since(start)
			respData := responseData{
				WithBody:          options.withBodies,
				WithStatus:        options.withStatus,
				WithContentLength: options.withContentLengths,
				WithDuration:      options.withDuration,
				Body:              string(wrapped.body),
				Status:            wrapped.status,
				ContentLength:     int64(len(wrapped.body)),
				Duration:          duration,
			}
			writer.Reset()
			err = respTemplate.Execute(&writer, respData)
			if err != nil {
				err := fmt.Errorf("error executing response template: %w", err)
				options.logger.Printf("Error: %v\n", err)
				http.Error(wrapped, err.Error(), http.StatusInternalServerError)
				return
			}
			options.logger.Println(writer.String())
		})
	}
}
