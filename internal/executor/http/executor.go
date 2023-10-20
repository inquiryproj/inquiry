package http

import (
	"log/slog"
	"net/http"
	"os"
)

type options struct {
	HTTPClient Client
	Logger     *slog.Logger
}

func defaultOptions() *options {
	return &options{
		HTTPClient: http.DefaultClient,
		Logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})),
	}
}

// Opts is a function for setting options on a scenario.
type Opts func(*options)

// WithHTTPClient sets the HTTP client to use for the scenario.
func WithHTTPClient(client Client) Opts {
	return func(o *options) {
		o.HTTPClient = client
	}
}

// WithLogger sets the logger to use for the scenario.
func WithLogger(logger *slog.Logger) Opts {
	return func(o *options) {
		o.Logger = logger
	}
}

// NewExecutor creates a new HTTP test scenario executor.
func NewExecutor(scenario *Scenario, opts ...Opts) (*Executor, error) {
	o := defaultOptions()

	for _, opt := range opts {
		opt(o)
	}
	executor := &Executor{}
	executor.scenario = scenario
	executor.httpClient = o.HTTPClient
	executor.logger = o.Logger

	return executor, nil
}
