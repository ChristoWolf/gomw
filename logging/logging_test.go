// Package logging_test implements tests for the logging package.
package logging_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/christowolf/go-middleware/logging"
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
			sut := logging.Middleware(logger)(handlerStub)
			sut.ServeHTTP(rec, req)
			if rec.Code != row.status {
				t.Errorf("expected status code %d, got %d", row.status, rec.Code)
			}
			if rec.Body.String() != row.resp {
				t.Errorf("expected response %s, got %s", row.resp, rec.Body.String())
			}
			if !strings.Contains(writer.String(), row.method) {
				t.Errorf("expected log to contain method %s, got %s", row.method, writer.String())
			}
			if !strings.Contains(writer.String(), row.url) {
				t.Errorf("expected log to contain url %s, got %s", row.url, writer.String())
			}
			if !strings.Contains(writer.String(), strconv.Itoa(row.status)) {
				t.Errorf("expected log to contain status %d, got %s", row.status, writer.String())
			}
			if !strings.Contains(writer.String(), row.resp) {
				t.Errorf("expected log to contain response %s, got %s", row.resp, writer.String())
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
	sut := logging.Middleware(logger)(handlerStub)
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
