// Package logging_test implements tests for the logging package.
package logging_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/christowolf/go-middleware/v2/logging"
)

// TestMiddleware tests the logging middleware.
func TestMiddleware(t *testing.T) {
	data := []struct {
		method string
		url    string
		status int
		resp   string
	}{
		{"GET", "/test/endpoint/get", http.StatusOK, "OK"},
		{"POST", "/test/endpoint/post", http.StatusCreated, "Created"},
		{"PUT", "/test/endpoint/put", http.StatusAccepted, "OK"},
		{"DELETE", "/test/endpoint/delete", http.StatusOK, "OK"},
		{"OPTIONS", "/test/endpoint/notfound", http.StatusNotFound, "Not found"},
	}
	t.Parallel()
	for _, row := range data {
		row := row
		t.Run(row.method+" "+row.url, func(t *testing.T) {
			t.Parallel()
			handlerStub := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(row.status)
				fmt.Fprint(w, row.resp)
			})
			req := httptest.NewRequest(row.method, row.url, nil)
			rec := httptest.NewRecorder()
			writer := strings.Builder{}
			logger := log.New(&writer, "", log.LstdFlags)
			options := logging.NewOptions(
				logging.WithLogger(logger),
				logging.WithBodies(true),
				logging.WithStatus(true),
				logging.WithMethod(true),
				logging.WithUrl(true),
				logging.WithContentLengths(true),
				logging.WithDuration(true))
			sut := logging.Middleware(*options)(handlerStub)
			sut.ServeHTTP(rec, req)
			if rec.Body.String() != row.resp {
				t.Errorf("expected response %s, got %s", row.resp, rec.Body.String())
			}
			if !strings.Contains(writer.String(), row.resp) {
				t.Errorf("expected log to contain response %s, got %s", row.resp, writer.String())
			}
			if rec.Code != row.status {
				t.Errorf("expected status code %d, got %d", row.status, rec.Code)
			}
			if !strings.Contains(writer.String(), strconv.Itoa(row.status)) {
				t.Errorf("expected log to contain status %d, got %s", row.status, writer.String())
			}
			if !strings.Contains(writer.String(), row.method) {
				t.Errorf("expected log to contain method %s, got %s", row.method, writer.String())
			}
			if !strings.Contains(writer.String(), row.url) {
				t.Errorf("expected log to contain url %s, got %s", row.url, writer.String())
			}
			if !strings.Contains(writer.String(), strconv.Itoa(len(row.resp))) {
				t.Errorf("expected log to contain content length %d, got %s", len(row.resp), writer.String())
			}
			if !strings.Contains(writer.String(), "µs") {
				t.Errorf("expected log to contain duration in microseconds, got %s", writer.String())
			}
		})
	}
}

// TestMiddlewareReadErr tests the
// logging middleware when the response body
// cannot be read.
func TestMiddlewareReadErr(t *testing.T) {
	t.Parallel()
	// We inject a handler which acts like everything worked.
	handlerStub := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
	req := httptest.NewRequest("GET", "/test/endpoint/error", errReader{})
	rec := httptest.NewRecorder()
	writer := strings.Builder{}
	logger := log.New(&writer, "", log.LstdFlags)
	options := logging.NewOptions(
		logging.WithLogger(logger),
		logging.WithBodies(true),
		logging.WithStatus(true),
		logging.WithMethod(true),
		logging.WithUrl(true),
		logging.WithContentLengths(true),
		logging.WithDuration(true))
	sut := logging.Middleware(*options)(handlerStub)
	sut.ServeHTTP(rec, req)
	wantCode := http.StatusInternalServerError
	if rec.Code != wantCode {
		t.Errorf("expected status code %d, got %d", wantCode, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "error") {
		t.Errorf("expected response body to contain error, got %s", rec.Body.String())
	}
	if !strings.Contains(writer.String(), "error") {
		t.Errorf("expected log to contain error, got %s", writer.String())
	}
}

// errReader is an io.Reader that always returns an error.
type errReader struct{}

// Read always returns an error.
func (errReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("refused to read")
}

// TestMiddlewareDuration tests if durations are
// logged with microseconds precision.
func TestMiddlewareDuration(t *testing.T) {
	data := []time.Duration{
		time.Nanosecond,
		time.Microsecond,
		time.Millisecond,
		time.Second,
	}
	for _, row := range data {
		row := row
		t.Run(row.String(), func(t *testing.T) {
			// These are not parallelized on purpose
			// as that would mess with the times
			// due to goroutine pausing.
			handlerStub := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(row)
			})
			req := httptest.NewRequest("GET", "/test/endpoint/duration", nil)
			rec := httptest.NewRecorder()
			writer := strings.Builder{}
			logger := log.New(&writer, "", log.LstdFlags)
			options := logging.NewOptions(
				logging.WithLogger(logger),
				logging.WithDuration(true))
			sut := logging.Middleware(*options)(handlerStub)
			sut.ServeHTTP(rec, req)
			logged := writer.String()
			regex := regexp.MustCompile(`\d+µs`)
			dur := regex.FindString(logged)
			if dur == "" {
				t.Errorf("expected duration to be logged, got %s", logged)
			}
			micro, err := time.ParseDuration(dur)
			if err != nil {
				t.Errorf("expected duration to be parsed, got error: %v", err)
			}
			if micro < row {
				t.Errorf("expected duration to be at least %s, got %s", row, micro)
			}
		})
	}
}
