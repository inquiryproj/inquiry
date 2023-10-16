// Package options provides the service options.
package options

import (
	"log/slog"
	"os"
)

// Options represents the options for the handlers.
type Options struct {
	Logger *slog.Logger
}

// DefaultOptions returns the default options.
func DefaultOptions() *Options {
	return &Options{
		Logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})),
	}
}

// Opts represents a function that modifies the options.
type Opts func(*Options)
