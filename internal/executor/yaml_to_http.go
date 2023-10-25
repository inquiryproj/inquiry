package executor

import (
	"github.com/inquiryproj/inquiry/internal/executor/http"
	"github.com/inquiryproj/inquiry/internal/executor/yaml"
)

func yamlScenarioToHTTPScenario(name string, yamlScenario *yaml.Scenario) *http.Scenario {
	return &http.Scenario{
		Name:  name,
		Steps: yamlStepsToHTTPSteps(yamlScenario.Steps),
	}
}

func yamlStepsToHTTPSteps(yamlSteps []*yaml.Step) []*http.Step {
	steps := []*http.Step{}
	for _, s := range yamlSteps {
		steps = append(steps, &http.Step{
			Name:       s.Name,
			Request:    yamlRequestToHTTPRequest(s.Request),
			Validation: yamlValidationToHTTPValidation(s.Validation),
			Retry:      yamlRetryToHTTPRetry(s.Retry),
		})
	}
	return steps
}

func yamlRetryToHTTPRetry(yamlRetry *yaml.Retry) *http.Retry {
	if yamlRetry == nil {
		return nil
	}
	return &http.Retry{
		Attempts: yamlRetry.Attempts,
		Timeout:  yamlRetry.Timeout,
	}
}

func yamlRequestToHTTPRequest(yamlRequest *yaml.Request) *http.Request {
	if yamlRequest == nil {
		return nil
	}
	return &http.Request{
		Method:  yamlRequest.Method,
		URL:     yamlRequest.URL,
		Headers: yamlHeadersToHTTPHeaders(yamlRequest.Headers),
		Body:    yamlRequest.Body,
	}
}

func yamlHeadersToHTTPHeaders(yamlHeaders []*yaml.Header) []*http.Header {
	headers := []*http.Header{}
	for _, h := range yamlHeaders {
		headers = append(headers, &http.Header{
			Name:  h.Name,
			Value: h.Value,
		})
	}
	return headers
}

func yamlValidationToHTTPValidation(yamlValidation *yaml.Validation) *http.Validation {
	if yamlValidation == nil {
		return nil
	}
	return &http.Validation{
		Body:    yamlAssertionsToHTTPAssertions(yamlValidation.Body),
		Status:  yamlAssertionToHTTPAssertion(yamlValidation.Status),
		Headers: yamlAssertionsToHTTPAssertions(yamlValidation.Headers),
	}
}

func yamlAssertionsToHTTPAssertions(yamlAssertions []*yaml.Assertion) []*http.Assertion {
	assertions := []*http.Assertion{}
	for _, a := range yamlAssertions {
		assertions = append(assertions, yamlAssertionToHTTPAssertion(a))
	}
	return assertions
}

func yamlAssertionToHTTPAssertion(yamlAssertion *yaml.Assertion) *http.Assertion {
	if yamlAssertion == nil {
		return nil
	}
	return &http.Assertion{
		Key:       yamlAssertion.Key,
		Assertion: http.AssertionMethod(string(yamlAssertion.Assertion)),
		Value:     yamlAssertion.Value,
	}
}
