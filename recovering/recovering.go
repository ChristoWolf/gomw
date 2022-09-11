// Package recovering implements a middleware
// for recovering from panics.
package recovering

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
)

const (
	// stackTraceWrapper is a text wrapper
	// around the debug stack for beautifying.
	stackTraceWrapper = `
	# Debug stack
	%s`
)

// Middleware is a middleware for recovering from panics.
func Middleware(options Options) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				var err error
				// According to recover's doc the
				// recovery was successful if
				// a non-nil error is returned.
				if rec := recover(); rec != nil {
					switch r := rec.(type) {
					case string:
						err = errors.New(r)
					case error:
						err = r
					default:
						err = errors.New("unknown recovery type")
					}
					err = fmt.Errorf("recovered from panic: %v", err)
					if options.withStackTrace {
						trace := fmt.Sprintf(stackTraceWrapper, debug.Stack())
						err = fmt.Errorf("%v%s", err, trace)
					}
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
