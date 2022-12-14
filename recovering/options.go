package recovering

// Options provides means to configure the middleware
// by applying the functional options pattern.
type Options struct {
	withStackTrace bool
}

// Option is a function that configures the middleware
// via the functional options pattern.
type Option func(*Options)

// NewOptions creates functionally injectable Options.
func NewOptions(options ...Option) *Options {
	opts := &Options{}
	for _, option := range options {
		option(opts)
	}
	return opts
}

// WithStackTrace configures whether to
// dump the stack trace after recovering.
// See the GOTRACEBACK environment variable
// (https://pkg.go.dev/runtime#hdr-Environment_Variables)
// on how to configure the stack trace level.
func WithStackTrace() Option {
	return func(o *Options) {
		o.withStackTrace = true
	}
}
