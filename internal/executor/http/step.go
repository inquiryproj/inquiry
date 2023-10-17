package http

import (
	"fmt"
	"io"
	"net/http"

	"github.com/tidwall/gjson"
)

func (s *Step) executeRequest(httpClient Client) (*RequestResult, error) {
	req, err := s.Request.toHTTPRequest()
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	s.IsExecuted = true
	s.RequestResult = &RequestResult{
		Body:    b,
		Status:  resp.StatusCode,
		Headers: resp.Header,
	}
	return s.RequestResult, nil
}

func (s Step) validate(requestResult *RequestResult) error {
	err := s.validateStatus(requestResult.Status)
	if err != nil {
		return err
	}

	err = s.validateHeaders(requestResult.Headers)
	if err != nil {
		return err
	}

	err = s.validateBody(requestResult.Body)
	if err != nil {
		return err
	}

	return nil
}

func (s Step) validateBody(body []byte) error {
	for _, assertion := range s.Validation.Body {
		jsonValue := gjson.Get(string(body), assertion.Key)
		if !jsonValue.Exists() {
			return s.errorForMsg(fmt.Sprintf("body key %s not found", assertion.Key))
		}

		err := s.assertValue(jsonValue.String(), ValidationBody, assertion)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s Step) validateStatus(status int) error {
	if s.Validation.Status == nil {
		return nil
	}

	return s.assertValue(fmt.Sprintf("%d", status), ValidationStatus, s.Validation.Status)
}

func (s Step) validateHeaders(headers http.Header) error {
	for _, assertion := range s.Validation.Headers {
		err := s.assertValue(headers.Get(assertion.Key), ValidationHeaders, assertion)
		if err != nil {
			return err
		}
	}
	return nil
}

// AssertionError is an error for when an assertion fails.
type AssertionError struct {
	StepName   string
	RequestURL string
	Msg        string
}

func (e AssertionError) Error() string {
	return fmt.Sprintf("failed at step \"%s\" for request to %s: %s", e.StepName, e.RequestURL, e.Msg)
}

func (s Step) assertValue(value string, validationType validationType, assertion *Assertion) error {
	switch assertion.Assertion {
	case AssertionMethodEqual:
		return s.assertEqual(value, validationType, assertion)
	case AssertionMethodRegex:
		// FIXME: implement regex assertions
	case AssertionMethodNotEmpty:
		return s.assertNotEmpty(value, validationType, assertion)
	}

	return nil
}

func (s Step) assertEqual(value string, validationType validationType, assertion *Assertion) error {
	if value == assertion.Value {
		return nil
	}
	switch validationType {
	case ValidationBody:
		return s.errorForMsg(fmt.Sprintf("body key %s has value %s, expected %s", assertion.Key, value, assertion.Value))
	case ValidationStatus:
		return s.errorForMsg(fmt.Sprintf("status has value %s, expected %s", value, assertion.Value))
	case ValidationHeaders:
		return s.errorForMsg(fmt.Sprintf("header %s has value %s, expected %s", assertion.Key, value, assertion.Value))
	}
	return nil
}

func (s Step) assertNotEmpty(value string, validationType validationType, assertion *Assertion) error {
	if value != "" {
		return nil
	}
	switch validationType {
	case ValidationBody:
		return s.errorForMsg(fmt.Sprintf("body key %s is empty", assertion.Key))
	case ValidationStatus:
		return s.errorForMsg("status has no value")
	case ValidationHeaders:
		return s.errorForMsg(fmt.Sprintf("header %s is empty", assertion.Key))
	}
	return nil
}

func (s Step) errorForMsg(msg string) error {
	return AssertionError{
		StepName:   s.Name,
		RequestURL: s.Request.URL,
		Msg:        msg,
	}
}
