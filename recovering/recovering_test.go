// Package recovering_test implements tests for the recovering package.
package recovering_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/christowolf/go-middleware/v2/recovering"
)

// TestMiddleware tests the recovering middleware
// by simulating a panic.
func TestMiddleware(t *testing.T) {
	t.Parallel()
	data := []any{
		"test panic",
		errors.New("test panic"),
		123,
	}
	for _, row := range data {
		row := row
		t.Run(fmt.Sprintf("%T", row), func(t *testing.T) {
			t.Parallel()
			// We inject a handler which panics.
			handlerStub := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic(row)
			})
			req := httptest.NewRequest("GET", "/test/endpoint/panic", nil)
			rec := httptest.NewRecorder()
			options := recovering.NewOptions(recovering.WithStackTrace())
			sut := recovering.Middleware(*options)(handlerStub)
			sut.ServeHTTP(rec, req)
			wantCode := http.StatusInternalServerError
			if rec.Code != wantCode {
				t.Errorf("expected status code %d, got %d", wantCode, rec.Code)
			}
			if !strings.Contains(rec.Body.String(), "panic") {
				t.Errorf("expected response body to contain panic, got %s", rec.Body.String())
			}
			if !strings.Contains(rec.Body.String(), "stack") {
				t.Errorf("expected response body to contain stack, got %s", rec.Body.String())
			}
		})
	}
}

// TestMiddlewareNoPanic tests the recovering middleware
// without any panic.
func TestMiddlewareNoPanic(t *testing.T) {
	t.Parallel()
	wantCode := http.StatusOK
	// We inject a handler which does not panic.
	handlerStub := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(wantCode)
	})
	req := httptest.NewRequest("GET", "/test/endpoint/mopanic", nil)
	rec := httptest.NewRecorder()
	options := recovering.NewOptions(recovering.WithStackTrace())
	sut := recovering.Middleware(*options)(handlerStub)
	sut.ServeHTTP(rec, req)
	if rec.Code != wantCode {
		t.Errorf("expected status code %d, got %d", wantCode, rec.Code)
	}
	if strings.Contains(rec.Body.String(), "panic") {
		t.Errorf("expected response body to not contain panic, got %s", rec.Body.String())
	}
}
