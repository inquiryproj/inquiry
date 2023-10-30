// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.15.0 DO NOT EDIT.
package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime"
)

const (
	ApiKeyAuthScopes = "ApiKeyAuth.Scopes"
)

// Defines values for ProjectRunOutputState.
const (
	Cancelled ProjectRunOutputState = "cancelled"
	Completed ProjectRunOutputState = "completed"
	Failure   ProjectRunOutputState = "failure"
	Pending   ProjectRunOutputState = "pending"
	Running   ProjectRunOutputState = "running"
)

// Defines values for ScenarioSpecType.
const (
	Yaml ScenarioSpecType = "yaml"
)

// ErrMsg defines model for ErrMsg.
type ErrMsg struct {
	Message string `json:"message"`
}

// Project defines model for Project.
type Project struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// ProjectArray defines model for ProjectArray.
type ProjectArray = []Project

// ProjectRunOutput defines model for ProjectRunOutput.
type ProjectRunOutput struct {
	ID                 uuid.UUID             `json:"id"`
	ProjectID          uuid.UUID             `json:"project_id"`
	ScenarioRunDetails []ScenarioRunDetails  `json:"scenario_run_details"`
	State              ProjectRunOutputState `json:"state"`
	Success            bool                  `json:"success"`
}

// ProjectRunOutputState defines model for ProjectRunOutput.State.
type ProjectRunOutputState string

// ProjectRunOutputArray defines model for ProjectRunOutputArray.
type ProjectRunOutputArray = []ProjectRunOutput

// ProjectRunRequest defines model for ProjectRunRequest.
type ProjectRunRequest struct {
	ProjectID   *uuid.UUID `json:"project_id,omitempty"`
	ProjectName *string    `json:"project_name,omitempty"`
}

// Scenario defines model for Scenario.
type Scenario struct {
	ID        uuid.UUID        `json:"id"`
	Name      string           `json:"name"`
	ProjectID uuid.UUID        `json:"project_id"`
	Spec      string           `json:"spec"`
	SpecType  ScenarioSpecType `json:"spec_type"`
}

// ScenarioSpecType defines model for Scenario.SpecType.
type ScenarioSpecType string

// ScenarioCreateRequest defines model for ScenarioCreateRequest.
type ScenarioCreateRequest struct {
	Name string `json:"name"`

	// Spec A base64 encoded string of the spec
	Spec     string `json:"spec"`
	SpecType string `json:"spec_type"`
}

// ScenarioRunDetails defines model for ScenarioRunDetails.
type ScenarioRunDetails struct {
	Assertions   int              `json:"assertions"`
	DurationInMs int              `json:"duration_in_ms"`
	Name         string           `json:"name"`
	Steps        []StepRunDetails `json:"steps"`
	Success      bool             `json:"success"`
}

// StepRunDetails defines model for StepRunDetails.
type StepRunDetails struct {
	Assertions          int    `json:"assertions"`
	DurationInMs        int    `json:"duration_in_ms"`
	Name                string `json:"name"`
	RequestDurationInMs int    `json:"request_duration_in_ms"`
	Retries             int    `json:"retries"`
	Success             bool   `json:"success"`
	URL                 string `json:"url"`
}

// ListProjectsParams defines parameters for ListProjects.
type ListProjectsParams struct {
	// Limit The number of projects to return
	Limit *int `form:"limit,omitempty" json:"limit,omitempty"`

	// Offset The number of projects to skip
	Offset *int `form:"offset,omitempty" json:"offset,omitempty"`
}

// ListRunsForProjectParams defines parameters for ListRunsForProject.
type ListRunsForProjectParams struct {
	// Limit The number of runs to return
	Limit *int `form:"limit,omitempty" json:"limit,omitempty"`

	// Offset The number of runs to skip
	Offset *int `form:"offset,omitempty" json:"offset,omitempty"`
}

// CreateProjectJSONRequestBody defines body for CreateProject for application/json ContentType.
type CreateProjectJSONRequestBody = Project

// RunProjectJSONRequestBody defines body for RunProject for application/json ContentType.
type RunProjectJSONRequestBody = ProjectRunRequest

// CreateScenarioJSONRequestBody defines body for CreateScenario for application/json ContentType.
type CreateScenarioJSONRequestBody = ScenarioCreateRequest

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// ListProjects request
	ListProjects(ctx context.Context, params *ListProjectsParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateProjectWithBody request with any body
	CreateProjectWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateProject(ctx context.Context, body CreateProjectJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RunProjectWithBody request with any body
	RunProjectWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	RunProject(ctx context.Context, body RunProjectJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// ListRunsForProject request
	ListRunsForProject(ctx context.Context, id uuid.UUID, params *ListRunsForProjectParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateScenarioWithBody request with any body
	CreateScenarioWithBody(ctx context.Context, id uuid.UUID, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateScenario(ctx context.Context, id uuid.UUID, body CreateScenarioJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) ListProjects(ctx context.Context, params *ListProjectsParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewListProjectsRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateProjectWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateProjectRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateProject(ctx context.Context, body CreateProjectJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateProjectRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RunProjectWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRunProjectRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RunProject(ctx context.Context, body RunProjectJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRunProjectRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) ListRunsForProject(ctx context.Context, id uuid.UUID, params *ListRunsForProjectParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewListRunsForProjectRequest(c.Server, id, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateScenarioWithBody(ctx context.Context, id uuid.UUID, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateScenarioRequestWithBody(c.Server, id, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateScenario(ctx context.Context, id uuid.UUID, body CreateScenarioJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateScenarioRequest(c.Server, id, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewListProjectsRequest generates requests for ListProjects
func NewListProjectsRequest(server string, params *ListProjectsParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/projects")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.Limit != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "limit", runtime.ParamLocationQuery, *params.Limit); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Offset != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "offset", runtime.ParamLocationQuery, *params.Offset); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewCreateProjectRequest calls the generic CreateProject builder with application/json body
func NewCreateProjectRequest(server string, body CreateProjectJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateProjectRequestWithBody(server, "application/json", bodyReader)
}

// NewCreateProjectRequestWithBody generates requests for CreateProject with any type of body
func NewCreateProjectRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/projects")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewRunProjectRequest calls the generic RunProject builder with application/json body
func NewRunProjectRequest(server string, body RunProjectJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewRunProjectRequestWithBody(server, "application/json", bodyReader)
}

// NewRunProjectRequestWithBody generates requests for RunProject with any type of body
func NewRunProjectRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/projects/run")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewListRunsForProjectRequest generates requests for ListRunsForProject
func NewListRunsForProjectRequest(server string, id uuid.UUID, params *ListRunsForProjectParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "id", runtime.ParamLocationPath, id)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/projects/%s/runs", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.Limit != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "limit", runtime.ParamLocationQuery, *params.Limit); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Offset != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "offset", runtime.ParamLocationQuery, *params.Offset); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewCreateScenarioRequest calls the generic CreateScenario builder with application/json body
func NewCreateScenarioRequest(server string, id uuid.UUID, body CreateScenarioJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateScenarioRequestWithBody(server, id, "application/json", bodyReader)
}

// NewCreateScenarioRequestWithBody generates requests for CreateScenario with any type of body
func NewCreateScenarioRequestWithBody(server string, id uuid.UUID, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "id", runtime.ParamLocationPath, id)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/projects/%s/scenarios", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// ListProjectsWithResponse request
	ListProjectsWithResponse(ctx context.Context, params *ListProjectsParams, reqEditors ...RequestEditorFn) (*ListProjectsResponse, error)

	// CreateProjectWithBodyWithResponse request with any body
	CreateProjectWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateProjectResponse, error)

	CreateProjectWithResponse(ctx context.Context, body CreateProjectJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateProjectResponse, error)

	// RunProjectWithBodyWithResponse request with any body
	RunProjectWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*RunProjectResponse, error)

	RunProjectWithResponse(ctx context.Context, body RunProjectJSONRequestBody, reqEditors ...RequestEditorFn) (*RunProjectResponse, error)

	// ListRunsForProjectWithResponse request
	ListRunsForProjectWithResponse(ctx context.Context, id uuid.UUID, params *ListRunsForProjectParams, reqEditors ...RequestEditorFn) (*ListRunsForProjectResponse, error)

	// CreateScenarioWithBodyWithResponse request with any body
	CreateScenarioWithBodyWithResponse(ctx context.Context, id uuid.UUID, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateScenarioResponse, error)

	CreateScenarioWithResponse(ctx context.Context, id uuid.UUID, body CreateScenarioJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateScenarioResponse, error)
}

type ListProjectsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *ProjectArray
	JSONDefault  *ErrMsg
}

// Status returns HTTPResponse.Status
func (r ListProjectsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ListProjectsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateProjectResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *Project
	JSONDefault  *ErrMsg
}

// Status returns HTTPResponse.Status
func (r CreateProjectResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateProjectResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RunProjectResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *ProjectRunOutput
	JSONDefault  *ErrMsg
}

// Status returns HTTPResponse.Status
func (r RunProjectResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RunProjectResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ListRunsForProjectResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *ProjectRunOutputArray
	JSONDefault  *ErrMsg
}

// Status returns HTTPResponse.Status
func (r ListRunsForProjectResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ListRunsForProjectResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateScenarioResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *Scenario
	JSONDefault  *ErrMsg
}

// Status returns HTTPResponse.Status
func (r CreateScenarioResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateScenarioResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// ListProjectsWithResponse request returning *ListProjectsResponse
func (c *ClientWithResponses) ListProjectsWithResponse(ctx context.Context, params *ListProjectsParams, reqEditors ...RequestEditorFn) (*ListProjectsResponse, error) {
	rsp, err := c.ListProjects(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseListProjectsResponse(rsp)
}

// CreateProjectWithBodyWithResponse request with arbitrary body returning *CreateProjectResponse
func (c *ClientWithResponses) CreateProjectWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateProjectResponse, error) {
	rsp, err := c.CreateProjectWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateProjectResponse(rsp)
}

func (c *ClientWithResponses) CreateProjectWithResponse(ctx context.Context, body CreateProjectJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateProjectResponse, error) {
	rsp, err := c.CreateProject(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateProjectResponse(rsp)
}

// RunProjectWithBodyWithResponse request with arbitrary body returning *RunProjectResponse
func (c *ClientWithResponses) RunProjectWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*RunProjectResponse, error) {
	rsp, err := c.RunProjectWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRunProjectResponse(rsp)
}

func (c *ClientWithResponses) RunProjectWithResponse(ctx context.Context, body RunProjectJSONRequestBody, reqEditors ...RequestEditorFn) (*RunProjectResponse, error) {
	rsp, err := c.RunProject(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRunProjectResponse(rsp)
}

// ListRunsForProjectWithResponse request returning *ListRunsForProjectResponse
func (c *ClientWithResponses) ListRunsForProjectWithResponse(ctx context.Context, id uuid.UUID, params *ListRunsForProjectParams, reqEditors ...RequestEditorFn) (*ListRunsForProjectResponse, error) {
	rsp, err := c.ListRunsForProject(ctx, id, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseListRunsForProjectResponse(rsp)
}

// CreateScenarioWithBodyWithResponse request with arbitrary body returning *CreateScenarioResponse
func (c *ClientWithResponses) CreateScenarioWithBodyWithResponse(ctx context.Context, id uuid.UUID, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateScenarioResponse, error) {
	rsp, err := c.CreateScenarioWithBody(ctx, id, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateScenarioResponse(rsp)
}

func (c *ClientWithResponses) CreateScenarioWithResponse(ctx context.Context, id uuid.UUID, body CreateScenarioJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateScenarioResponse, error) {
	rsp, err := c.CreateScenario(ctx, id, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateScenarioResponse(rsp)
}

// ParseListProjectsResponse parses an HTTP response from a ListProjectsWithResponse call
func ParseListProjectsResponse(rsp *http.Response) (*ListProjectsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &ListProjectsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest ProjectArray
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest ErrMsg
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseCreateProjectResponse parses an HTTP response from a CreateProjectWithResponse call
func ParseCreateProjectResponse(rsp *http.Response) (*CreateProjectResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateProjectResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest Project
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest ErrMsg
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseRunProjectResponse parses an HTTP response from a RunProjectWithResponse call
func ParseRunProjectResponse(rsp *http.Response) (*RunProjectResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RunProjectResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest ProjectRunOutput
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest ErrMsg
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseListRunsForProjectResponse parses an HTTP response from a ListRunsForProjectWithResponse call
func ParseListRunsForProjectResponse(rsp *http.Response) (*ListRunsForProjectResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &ListRunsForProjectResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest ProjectRunOutputArray
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest ErrMsg
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseCreateScenarioResponse parses an HTTP response from a CreateScenarioWithResponse call
func ParseCreateScenarioResponse(rsp *http.Response) (*CreateScenarioResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateScenarioResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest Scenario
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest ErrMsg
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}
