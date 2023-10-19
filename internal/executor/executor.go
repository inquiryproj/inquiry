// Package executor contains the test executor app.
package executor

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/inquiryproj/inquiry/internal/executor/http"
	"github.com/inquiryproj/inquiry/internal/infra/replacer"
	"github.com/inquiryproj/inquiry/internal/infra/test/definitions/yaml"
)

// error definitions.
var (
	ErrCreateExecutor = fmt.Errorf("unable to create test executor")
)

// App is the interface for the test executor app.
type App interface {
	Play() (*http.ExecuteResult, error)
}

type options struct {
	Reader io.Reader
	Logger *slog.Logger
}

func defaultOptions() *options {
	return &options{
		Reader: os.Stdin,
		Logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})),
	}
}

// Opts is a function for setting options on a test exectuor.
type Opts func(*options)

// WithLogger sets the logger to use for the scenario.
func WithLogger(logger *slog.Logger) Opts {
	return func(o *options) {
		o.Logger = logger
	}
}

// WithReader sets the reader to use for the scenario.
func WithReader(reader io.Reader) Opts {
	return func(o *options) {
		o.Reader = reader
	}
}

// New creates a new test executor app.
func New(name string, opts ...Opts) (App, error) {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	data, err := io.ReadAll(o.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTTP scenario definition: %w", err)
	}
	testSpec, yamlScenario, err := readData(data)
	if err != nil {
		return nil, fmt.Errorf("failed to read test definition: %w", err)
	}
	return newAppForTestDefinition(name, testSpec, yamlScenario, o)
}

type scenarioExecutor interface {
	Play() (*http.ExecuteResult, error)
}

type app struct {
	scenarioExecutor scenarioExecutor
}

func (a *app) Play() (*http.ExecuteResult, error) {
	return a.scenarioExecutor.Play()
}

func newAppForTestDefinition(name string,
	testSpec *TestSpec, scenario *yaml.Scenario,
	options *options,
) (*app, error) {
	switch testSpec.Type {
	case TestTypeHTTP:
		httpExecutor, err := http.NewExecutor(
			yamlScenarioToHTTPScenario(name, scenario),
			http.WithLogger(options.Logger),
		)
		if err != nil {
			return nil, err
		}
		return &app{
			scenarioExecutor: httpExecutor,
		}, nil
	default:
		return nil, ErrCreateExecutor
	}
}

func readData(data []byte) (*TestSpec, *yaml.Scenario, error) {
	yamlTestSpec, yamlScenario, err := yaml.NewTestDefinitionFromBytes(
		data,
		replacer.NewFuncReplacer(),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse HTTP scenario definition: %w", err)
	}

	return yamlTestSpecToTestSpec(yamlTestSpec), yamlScenario, nil
}
