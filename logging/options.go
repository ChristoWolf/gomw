package logging

import (
	"log"
	"os"
)

// Options provides means to configure the middleware
// by applying the functional options pattern.
// TODO: Add headers support.
type Options struct {
	logger             *log.Logger
	withBodies         bool
	withStatus         bool
	withMethod         bool
	withUrl            bool
	withContentLengths bool
	withDuration       bool
}

// NewOptions creates functionally injectable Options.
// If no logger is provided, it defaults to os.Stdout.
func NewOptions(options ...func(*Options)) *Options {
	opts := &Options{logger: log.New(os.Stdout, "", log.LstdFlags)}
	for _, option := range options {
		option(opts)
	}
	return opts
}

// WithLogger configures the logger to use.
func WithLogger(logger *log.Logger) func(*Options) {
	return func(o *Options) {
		o.logger = logger
	}
}

// WithBodies configures whether to log request and response bodies.
func WithBodies(withBodies bool) func(*Options) {
	return func(o *Options) {
		o.withBodies = withBodies
	}
}

// WithStatus configures whether to log the response status.
func WithStatus(withStatus bool) func(*Options) {
	return func(o *Options) {
		o.withStatus = withStatus
	}
}

// WithMethod configures whether to log the request method.
func WithMethod(withMethod bool) func(*Options) {
	return func(o *Options) {
		o.withMethod = withMethod
	}
}

// WithUrl configures whether to log the request URL.
func WithUrl(withUrl bool) func(*Options) {
	return func(o *Options) {
		o.withUrl = withUrl
	}
}

// WithContentLengths configures whether to log
// the request's and response's content lengths.
func WithContentLengths(withContentLengths bool) func(*Options) {
	return func(o *Options) {
		o.withContentLengths = withContentLengths
	}
}

// WithDuration configures whether to log the request duration.
func WithDuration(withDuration bool) func(*Options) {
	return func(o *Options) {
		o.withDuration = withDuration
	}
}
