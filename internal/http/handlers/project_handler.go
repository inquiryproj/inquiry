package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/inquiryproj/inquiry/internal/app"
	httpInternal "github.com/inquiryproj/inquiry/internal/http"
	"github.com/inquiryproj/inquiry/internal/service"
)

// ProjectHandler handles project requests.
type ProjectHandler struct {
	projectService service.Project
	logger         *slog.Logger
}

// newProjectHandler creates a new project handler.
func newProjectHandler(projectService service.Project, options *Options) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
		logger:         options.Logger,
	}
}

// ListProjects lists all projects.
func (h *ProjectHandler) ListProjects(ctx echo.Context, params httpInternal.ListProjectsParams) error {
	getProjectsRequest := &app.GetProjectsRequest{
		Limit:  100,
		Offset: 0,
	}
	if params.Limit != nil {
		getProjectsRequest.Limit = *params.Limit
	}
	if params.Offset != nil {
		getProjectsRequest.Offset = *params.Offset
	}

	projects, err := h.projectService.GetProjects(ctx.Request().Context(), getProjectsRequest)
	if err != nil {
		h.logger.Error("unable to get projects", slog.String("error", err.Error()))
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to get projects")
	}

	result := make([]httpInternal.Project, len(projects))
	for i, project := range projects {
		result[i] = httpInternal.Project{
			ID:   project.ID,
			Name: project.Name,
		}
	}

	return ctx.JSON(http.StatusOK, result)
}

// CreateProject creates a new project.
func (h *ProjectHandler) CreateProject(ctx echo.Context) error {
	httpProject := &httpInternal.Project{}
	err := json.NewDecoder(ctx.Request().Body).Decode(&httpProject)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid create project payload")
	}
	project, err := h.projectService.CreateProject(ctx.Request().Context(), &app.CreateProjectRequest{
		Name: httpProject.Name,
	})
	if errors.Is(err, app.ErrProjectAlreadyExists) {
		return echo.NewHTTPError(http.StatusConflict, fmt.Sprintf("project with name %s already exists", httpProject.Name))
	}
	if err != nil {
		h.logger.Error("unable to create project", slog.String("error", err.Error()))
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to create project")
	}
	return ctx.JSON(http.StatusCreated, httpInternal.Project{
		ID:   project.ID,
		Name: project.Name,
	})
}
