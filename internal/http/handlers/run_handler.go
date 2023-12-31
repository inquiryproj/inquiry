package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/inquiryproj/inquiry/internal/app"
	"github.com/inquiryproj/inquiry/internal/http/api"
	"github.com/inquiryproj/inquiry/internal/service"
)

// RunHandler handles run requests.
type RunHandler struct {
	runnerService service.Runner

	logger *slog.Logger
}

// newRunHandler creates a new run handler.
func newRunHandler(runnerService service.Runner, opts ...Opts) *RunHandler {
	options := defaultOptions()
	for _, o := range opts {
		o(options)
	}
	return &RunHandler{
		runnerService: runnerService,
		logger:        options.Logger,
	}
}

// RunProject runs all scenarios for a given project.
func (h *RunHandler) RunProject(ctx echo.Context) error {
	runProjectJSONRequestBody := api.RunProjectJSONRequestBody{}
	err := json.NewDecoder(ctx.Request().Body).Decode(&runProjectJSONRequestBody)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid create scenario payload")
	}
	if runProjectJSONRequestBody.ProjectID == nil && runProjectJSONRequestBody.ProjectName == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "either project_id or project_name must be provided")
	}
	var projectRunOutput *app.ProjectRunOutput
	if runProjectJSONRequestBody.ProjectID != nil {
		projectRunOutput, err = h.runnerService.RunProject(ctx.Request().Context(), &app.RunProjectRequest{
			ProjectID: *runProjectJSONRequestBody.ProjectID,
		})
	} else {
		projectRunOutput, err = h.runnerService.RunProjectByName(ctx.Request().Context(), &app.RunProjectByNameRequest{
			ProjectName: *runProjectJSONRequestBody.ProjectName,
		})
	}
	if errors.Is(err, app.ErrProjectNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "project not found")
	} else if err != nil {
		h.logger.Error("failed to run project", slog.String("error", err.Error()))
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to run project")
	}
	return ctx.JSON(http.StatusOK, projectRunOutputToHTTP(projectRunOutput))
}

func projectRunOutputToHTTP(projectRunOutput *app.ProjectRunOutput) api.ProjectRunOutput {
	return api.ProjectRunOutput{
		ID:        projectRunOutput.ID,
		ProjectID: projectRunOutput.ProjectID,
		Success:   projectRunOutput.Success,
		State:     api.ProjectRunOutputState(projectRunOutput.State),
	}
}

// ListRunsForProject returns runs for a given project.
func (h *RunHandler) ListRunsForProject(ctx echo.Context, id uuid.UUID, params api.ListRunsForProjectParams) error {
	listRunsForProjectRequest := &app.ListRunsForProjectRequest{
		Limit:     100,
		Offset:    0,
		ProjectID: id,
	}
	if params.Limit != nil {
		listRunsForProjectRequest.Limit = *params.Limit
	}
	if params.Offset != nil {
		listRunsForProjectRequest.Offset = *params.Offset
	}
	runs, err := h.runnerService.ListRunsForProject(ctx.Request().Context(), listRunsForProjectRequest)
	if err != nil {
		h.logger.Error("failed to get runs for project", slog.String("error", err.Error()))
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to get runs for project")
	}

	result := make([]api.ProjectRunOutput, len(runs.Runs))
	for i, run := range runs.Runs {
		result[i] = api.ProjectRunOutput{
			ID:                 run.ID,
			ProjectID:          run.ProjectID,
			Success:            run.Success,
			State:              api.ProjectRunOutputState(run.State),
			ScenarioRunDetails: appScenarioDetailsToHTTPScenarioDetails(run.ScenarioRunDetails),
		}
	}

	return ctx.JSON(http.StatusOK, result)
}

func appScenarioDetailsToHTTPScenarioDetails(scenario []*app.ScenarioRunDetails) []api.ScenarioRunDetails {
	result := []api.ScenarioRunDetails{}
	for _, detail := range scenario {
		result = append(result, api.ScenarioRunDetails{
			Name:         detail.Name,
			DurationInMs: int(detail.Duration.Milliseconds()),
			Assertions:   detail.Assertions,
			Steps:        appStepsRunDetailsToHTTPStepRunDetails(detail.Steps),
			Success:      detail.Success,
		})
	}
	return result
}

func appStepsRunDetailsToHTTPStepRunDetails(steps []*app.StepRunDetails) []api.StepRunDetails {
	result := []api.StepRunDetails{}
	for _, detail := range steps {
		result = append(result, api.StepRunDetails{
			Name:                detail.Name,
			Assertions:          detail.Assertions,
			URL:                 detail.URL,
			RequestDurationInMs: int(detail.RequestDuration.Milliseconds()),
			DurationInMs:        int(detail.Duration.Milliseconds()),
			Retries:             detail.Retries,
			Success:             detail.Success,
		})
	}
	return result
}
