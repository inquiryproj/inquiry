// Package http contains the HTTP/REST test scenario implementation.
package http

import (
	"log/slog"
	"net/http"
	"time"
)

type assertionMethod string

// Different assertion methods.
const (
	AssertionMethodEqual    assertionMethod = "equal"
	AssertionMethodRegex    assertionMethod = "regex"
	AssertionMethodNotEmpty assertionMethod = "not_empty"
)

// AssertionMethod returns an assertion method for a given string.
func AssertionMethod(method string) assertionMethod { //nolint: revive
	return assertionMethod(method)
}

type validationType string

// Different validation types.
const (
	ValidationBody    validationType = "body"
	ValidationStatus  validationType = "status"
	ValidationHeaders validationType = "headers"
)

// Client is the interface for perfoming HTTP requests.
type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

// Executor is the http test executor implementation.
type Executor struct {
	scenario   *Scenario
	httpClient Client
	logger     *slog.Logger
}

// Scenario is the main struct for a test scenario to be executed.
type Scenario struct {
	Name  string
	Steps []*Step
}

// ScenarioMetrics is a struct for storing metrics of a scenario.
type scenarioMetrics struct {
	TotalExecutionTime time.Duration
}

// InputReplacement is a struct for replacing dynamic inputs in a step.
type InputReplacement struct {
	StepName         string
	JSONKey          string
	ReplacementValue string
}

// Step represents a step in a scenario.
type Step struct {
	Name          string
	Request       *Request
	Validation    *Validation
	RequestResult *RequestResult
	IsExecuted    bool
	Retry         *Retry
}

// Retry for a single step.
type Retry struct {
	Attempts int           `yaml:"attempts"`
	Timeout  time.Duration `yaml:"timeout"`
}

// Request represents a HTTP request.
type Request struct {
	Method  string
	URL     string
	Headers []*Header
	Body    string
}

// RequestResult represents the result of an HTTP request.
type RequestResult struct {
	Body    []byte
	Status  int
	Headers http.Header
}

// Validation represents the validation of a step.
type Validation struct {
	Body    []*Assertion
	Status  *Assertion
	Headers []*Assertion
}

// Assertion represents an assertion, as part of a validation.
type Assertion struct {
	Key       string
	Assertion assertionMethod
	Value     string
}

// Header represents a HTTP header.
type Header struct {
	Name  string
	Value string
}
