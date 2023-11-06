// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.15.0 DO NOT EDIT.
package api

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /v1/projects)
	ListProjects(ctx echo.Context, params ListProjectsParams) error

	// (POST /v1/projects)
	CreateProject(ctx echo.Context) error

	// (POST /v1/projects/run)
	RunProject(ctx echo.Context) error

	// (GET /v1/projects/{id}/runs)
	ListRunsForProject(ctx echo.Context, id uuid.UUID, params ListRunsForProjectParams) error

	// (GET /v1/projects/{project_id}/scenarios)
	ListScenariosForProject(ctx echo.Context, projectId uuid.UUID, params ListScenariosForProjectParams) error

	// (POST /v1/projects/{project_id}/scenarios)
	CreateScenario(ctx echo.Context, projectId uuid.UUID) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// ListProjects converts echo context to params.
func (w *ServerInterfaceWrapper) ListProjects(ctx echo.Context) error {
	var err error

	ctx.Set(ApiKeyAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params ListProjectsParams
	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	// ------------- Optional query parameter "offset" -------------

	err = runtime.BindQueryParameter("form", true, false, "offset", ctx.QueryParams(), &params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter offset: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ListProjects(ctx, params)
	return err
}

// CreateProject converts echo context to params.
func (w *ServerInterfaceWrapper) CreateProject(ctx echo.Context) error {
	var err error

	ctx.Set(ApiKeyAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.CreateProject(ctx)
	return err
}

// RunProject converts echo context to params.
func (w *ServerInterfaceWrapper) RunProject(ctx echo.Context) error {
	var err error

	ctx.Set(ApiKeyAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.RunProject(ctx)
	return err
}

// ListRunsForProject converts echo context to params.
func (w *ServerInterfaceWrapper) ListRunsForProject(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id uuid.UUID

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	ctx.Set(ApiKeyAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params ListRunsForProjectParams
	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	// ------------- Optional query parameter "offset" -------------

	err = runtime.BindQueryParameter("form", true, false, "offset", ctx.QueryParams(), &params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter offset: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ListRunsForProject(ctx, id, params)
	return err
}

// ListScenariosForProject converts echo context to params.
func (w *ServerInterfaceWrapper) ListScenariosForProject(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "project_id" -------------
	var projectId uuid.UUID

	err = runtime.BindStyledParameterWithLocation("simple", false, "project_id", runtime.ParamLocationPath, ctx.Param("project_id"), &projectId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter project_id: %s", err))
	}

	ctx.Set(ApiKeyAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params ListScenariosForProjectParams
	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	// ------------- Optional query parameter "offset" -------------

	err = runtime.BindQueryParameter("form", true, false, "offset", ctx.QueryParams(), &params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter offset: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ListScenariosForProject(ctx, projectId, params)
	return err
}

// CreateScenario converts echo context to params.
func (w *ServerInterfaceWrapper) CreateScenario(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "project_id" -------------
	var projectId uuid.UUID

	err = runtime.BindStyledParameterWithLocation("simple", false, "project_id", runtime.ParamLocationPath, ctx.Param("project_id"), &projectId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter project_id: %s", err))
	}

	ctx.Set(ApiKeyAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.CreateScenario(ctx, projectId)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/v1/projects", wrapper.ListProjects)
	router.POST(baseURL+"/v1/projects", wrapper.CreateProject)
	router.POST(baseURL+"/v1/projects/run", wrapper.RunProject)
	router.GET(baseURL+"/v1/projects/:id/runs", wrapper.ListRunsForProject)
	router.GET(baseURL+"/v1/projects/:project_id/scenarios", wrapper.ListScenariosForProject)
	router.POST(baseURL+"/v1/projects/:project_id/scenarios", wrapper.CreateScenario)

}
