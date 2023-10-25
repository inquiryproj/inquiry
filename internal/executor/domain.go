package executor

import (
	"github.com/inquiryproj/inquiry/internal/executor/yaml"
)

type testType string

// Different test types.
const (
	TestTypeHTTP testType = "http"
)

// TestSpec for a single test scenario.
type TestSpec struct {
	Version   string
	Type      testType
	Variables []*Variable
}

// Variable for a single test definition.
type Variable struct {
	Name  string
	Value string
}

func yamlTestSpecToTestSpec(testDefinition *yaml.TestSpec) *TestSpec {
	return &TestSpec{
		Version: testDefinition.Version,
		Type:    testType(testDefinition.Type),
		Variables: func() []*Variable {
			variables := []*Variable{}
			for _, v := range testDefinition.Variables {
				variables = append(variables, &Variable{
					Name:  v.Name,
					Value: v.Value,
				})
			}
			return variables
		}(),
	}
}
