// Package yaml presents the domain model for the YAML format representing test scenarios.
package yaml

import (
	"fmt"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/inquiryproj/inquiry/internal/executor/replacer"
)

const variablesPrefix = "variables"

type assertionMethod string

// Different assertion methods.
const (
	AssertionMethodEqual    assertionMethod = "equal"
	AssertionMethodRegex    assertionMethod = "regex"
	AssertionMethodNotEmpty assertionMethod = "not_empty"
)

type testType string

// Different test types.
const (
	TestTypeHTTP testType = "http"
)

// Scenario represents a single test scenario represented in YAML.
type Scenario struct {
	Steps []*Step `yaml:"steps"`
}

// TestSpec for a single scenario.
type TestSpec struct {
	Version   string      `yaml:"version"`
	Type      testType    `yaml:"type"`
	Variables []*Variable `yaml:"variables"`
}

func (v TestSpec) getVariablesMap() map[string]string {
	variables := make(map[string]string)
	for _, v := range v.Variables {
		variables[fmt.Sprintf("%s.%s", variablesPrefix, v.Name)] = v.Value
	}
	return variables
}

// Variable for a single scenario.
type Variable struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// Step for a single scenario.
type Step struct {
	Name       string      `yaml:"name"`
	Request    *Request    `yaml:"request"`
	Validation *Validation `yaml:"validation"`
	Retry      *Retry      `yaml:"retry"`
}

// Retry for a single step.
type Retry struct {
	Attempts int           `yaml:"attempts"`
	Timeout  time.Duration `yaml:"timeout"`
}

// Request for a single step.
type Request struct {
	Method  string    `yaml:"method"`
	URL     string    `yaml:"url"`
	Headers []*Header `yaml:"headers"`
	Body    string    `yaml:"body"`
}

// Header for a request.
type Header struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// Validation for a single step.
type Validation struct {
	Body    []*Assertion `yaml:"body"`
	Status  *Assertion   `yaml:"status"`
	Headers []*Assertion `yaml:"headers"`
}

// Assertion represents an assertion, as part of a validation.
type Assertion struct {
	Key       string          `yaml:"key"`
	Assertion assertionMethod `yaml:"assertion"`
	Value     string          `yaml:"value"`
}

// NewTestDefinitionFromBytes creates a new test definition from a byte array representing a
// YAML file.
func NewTestDefinitionFromBytes(data []byte, replacers ...replacer.Replacer) (*TestSpec, *Scenario, error) {
	var testSpec TestSpec

	err := yaml.Unmarshal(data, &testSpec)
	if err != nil {
		return nil, nil, err
	}
	fileContent := string(data)
	replacers = append(replacers, replacer.NewMapReplacer(testSpec.getVariablesMap()))
	for _, r := range replacers {
		fileContent = r.Replace(fileContent)
	}

	var scenario Scenario

	err = yaml.Unmarshal([]byte(fileContent), &scenario)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid yaml definition for scenario after parsing variables %w", err)
	}

	return &testSpec, &scenario, nil
}
