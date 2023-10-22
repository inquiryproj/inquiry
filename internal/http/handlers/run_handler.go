package handlers

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/inquiryproj/inquiry/internal/app"
	httpInternal "github.com/inquiryproj/inquiry/internal/http"
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
func (h *RunHandler) RunProject(ctx echo.Context, id uuid.UUID) error {
	projectRunOutput, err := h.runnerService.RunProject(ctx.Request().Context(), &app.RunProjectRequest{
		ProjectID: id,
	})
	if err != nil {
		h.logger.Error("failed to run project", slog.String("error", err.Error()))
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to run project")
	}
	return ctx.JSON(http.StatusOK, httpInternal.ProjectRunOutput{
		ID:        projectRunOutput.ID,
		ProjectID: projectRunOutput.ProjectID,
		Success:   projectRunOutput.Success,
		State:     httpInternal.ProjectRunOutputState(projectRunOutput.State),
	})
}

// GetRunsForProject returns runs for a given project.
func (h *RunHandler) GetRunsForProject(ctx echo.Context, id uuid.UUID, params httpInternal.GetRunsForProjectParams) error {
	getRunsForProjectRequest := &app.GetRunsForProjectRequest{
		Limit:     100,
		Offset:    0,
		ProjectID: id,
	}
	if params.Limit != nil {
		getRunsForProjectRequest.Limit = *params.Limit
	}
	if params.Offset != nil {
		getRunsForProjectRequest.Offset = *params.Offset
	}
	runs, err := h.runnerService.GetRunsForProject(ctx.Request().Context(), getRunsForProjectRequest)
	if err != nil {
		h.logger.Error("failed to get runs for project", slog.String("error", err.Error()))
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to get runs for project")
	}

	result := make([]httpInternal.ProjectRunOutput, len(runs.Runs))
	for i, run := range runs.Runs {
		result[i] = httpInternal.ProjectRunOutput{
			ID:                 run.ID,
			ProjectID:          run.ProjectID,
			Success:            run.Success,
			State:              httpInternal.ProjectRunOutputState(run.State),
			ScenarioRunDetails: appScenarioDetailsToHTTPScenarioDetails(run.ScenarioRunDetails),
		}
	}

	return ctx.JSON(http.StatusOK, result)
}

func appScenarioDetailsToHTTPScenarioDetails(scenario []*app.ScenarioRunDetails) []httpInternal.ScenarioRunDetails {
	result := []httpInternal.ScenarioRunDetails{}
	for _, detail := range scenario {
		result = append(result, httpInternal.ScenarioRunDetails{
			Name:         detail.Name,
			DurationInMs: int(detail.Duration.Milliseconds()),
			Assertions:   detail.Assertions,
			Steps:        appStepsRunDetailsToHTTPStepRunDetails(detail.Steps),
			Success:      detail.Success,
		})
	}
	return result
}

func appStepsRunDetailsToHTTPStepRunDetails(steps []*app.StepRunDetails) []httpInternal.StepRunDetails {
	result := []httpInternal.StepRunDetails{}
	for _, detail := range steps {
		result = append(result, httpInternal.StepRunDetails{
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
