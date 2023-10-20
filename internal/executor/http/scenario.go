package http

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

// ErrInvalidPathReplacement is an error for when a replacement json path
// is not found in the defined step's JSON output.
type ErrInvalidPathReplacement struct {
	path string
}

func (e ErrInvalidPathReplacement) Error() string {
	return fmt.Sprintf("invalid path %s", e.path)
}

// ErrNonExecutedStep is an error for when a step is not executed but is required
// for dynamic input replacement.
type ErrNonExecutedStep struct {
	StepName, NonExecutedStepName string
}

func (e ErrNonExecutedStep) Error() string {
	return fmt.Sprintf("step %s requires input from non executed step %s", e.StepName, e.NonExecutedStepName)
}

// ErrJSONKeyNotFound is an error for when a JSON key for dynamic input
// replacement is not found in the defined step's JSON output.
type ErrJSONKeyNotFound struct {
	StepName, Key, Body string
}

func (e ErrJSONKeyNotFound) Error() string {
	return fmt.Sprintf("key %s not found in %s for step %s", e.Key, e.Body, e.StepName)
}

// ExecuteResult is the result of executing a scenario.
type ExecuteResult struct {
	Name               string
	TotalExecutionTime time.Duration
	TotalAssertions    int
	StepResults        []*ExecuteStepResult
	Success            bool
}

// ExecuteStepResult is the result of executing a step.
type ExecuteStepResult struct {
	Name            string
	Assertions      int
	URL             string
	RequestDuration time.Duration
	Duration        time.Duration
	Retries         int
	Success         bool
}

// Play executes the scenario.
func (e Executor) Play() (*ExecuteResult, error) {
	executeResult := &ExecuteResult{
		Name: e.scenario.Name,
	}
	start := time.Now()
	totalAssertions := 0
	success := true
	for _, step := range e.scenario.Steps {
		stepResult, err := e.playStep(step)
		totalAssertions += stepResult.Assertions
		success = stepResult.Success && success
		executeResult.StepResults = append(executeResult.StepResults, stepResult)
		if err != nil {
			break
		}
	}
	executeResult.TotalExecutionTime = time.Since(start)
	executeResult.TotalAssertions = totalAssertions
	executeResult.Success = success
	return executeResult, nil
}

// FIXME don't return error on validation failures, distinct in stepresult.
func (e Executor) playStep(step *Step) (*ExecuteStepResult, error) {
	err := e.replaceDynamicInputs(step)
	if err != nil {
		return &ExecuteStepResult{}, err
	}

	retries := 0
	if step.Retry != nil {
		retries = step.Retry.Attempts
	}
	timeout := 1 * time.Second
	if step.Retry != nil {
		timeout = step.Retry.Timeout
	}
	start := time.Now()
	stepResult, err := e.executeWithRetries(retries, timeout, step)

	stepResult.Duration = time.Since(start)
	return stepResult, err
}

func (e Executor) executeWithRetries(retries int, timeout time.Duration, step *Step) (*ExecuteStepResult, error) {
	retriesLeft := 0
	if step.Retry != nil {
		retriesLeft = step.Retry.Attempts - retries
	}

	stepResult, err := e.executeAndValidate(step)
	stepResult.Retries = retriesLeft
	if retries <= 0 {
		return stepResult, err
	}
	if err != nil {
		e.logger.Debug(fmt.Sprintf("retrying step %s in %v seconds", step.Name, timeout.Seconds()))
		time.Sleep(timeout)
		return e.executeWithRetries(retries-1, timeout, step)
	}

	return stepResult, nil
}

func (e Executor) executeAndValidate(step *Step) (*ExecuteStepResult, error) {
	stepResult := &ExecuteStepResult{
		Name:       step.Name,
		URL:        step.Request.URL,
		Assertions: len(step.Validation.Body) + len(step.Validation.Headers) + 1,
		Success:    false,
	}
	start := time.Now()
	requestResult, err := step.executeRequest(e.httpClient)
	stepResult.RequestDuration = time.Since(start)
	if err != nil {
		return stepResult, err
	}

	err = step.validate(requestResult)
	if err == nil {
		stepResult.Success = true
	}
	return stepResult, nil
}

func (e Executor) replaceDynamicInputs(step *Step) error {
	b, err := json.Marshal(step)
	if err != nil {
		return err
	}
	stepJSONString := string(b)
	replaceKeyMap, err := e.createReplacementMap(stepJSONString)
	if err != nil {
		return err
	}
	err = e.scenario.findReplacementValues(step.Name, replaceKeyMap)
	if err != nil {
		return err
	}

	for k, v := range replaceKeyMap {
		stepJSONString = strings.ReplaceAll(stepJSONString, k, v.ReplacementValue)
	}
	return json.Unmarshal([]byte(stepJSONString), &step)
}

func (e Executor) createReplacementMap(stepJSONString string) (map[string]*InputReplacement, error) {
	dynamicInputPlaceHolders := regexp.MustCompile(`\$\{steps.([^\}]*)\}`).FindAllStringSubmatch(stepJSONString, -1)
	replaceKeyMap := map[string]*InputReplacement{}
	for _, placeHolder := range dynamicInputPlaceHolders {
		if len(placeHolder) != 2 {
			e.logger.Warn("invalid dynamic placeholder detected", slog.String("placeholder", strings.Join(placeHolder, ",")))
			continue
		}
		separatedPath := strings.Split(placeHolder[1], ".")
		if len(separatedPath) <= 3 {
			return nil, ErrInvalidPathReplacement{
				path: placeHolder[1],
			}
		}
		replaceKeyMap[placeHolder[0]] = &InputReplacement{
			StepName: separatedPath[0],
			JSONKey:  strings.Join(separatedPath[3:], "."),
		}
	}

	return replaceKeyMap, nil
}

func (s Scenario) findReplacementValues(stepName string, replaceKeyMap map[string]*InputReplacement) error {
	for _, v := range replaceKeyMap {
		for _, s := range s.Steps {
			if s.Name == v.StepName {
				if !s.IsExecuted {
					return ErrNonExecutedStep{
						StepName:            stepName,
						NonExecutedStepName: v.StepName,
					}
				}
				jsonValue := gjson.Get(string(s.RequestResult.Body), v.JSONKey)
				if !jsonValue.Exists() {
					return ErrJSONKeyNotFound{
						StepName: stepName,
						Key:      v.JSONKey,
						Body:     string(s.RequestResult.Body),
					}
				}
				v.ReplacementValue = jsonValue.String()
			}
		}
	}
	return nil
}
