package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/inquiryproj/inquiry/internal/app"
	"github.com/inquiryproj/inquiry/internal/http/api"
	"github.com/inquiryproj/inquiry/internal/service"
)

// ScenarioHandler handles scenario requests.
type ScenarioHandler struct {
	scenarioService service.Scenario
	logger          *slog.Logger
}

// newScenarioHandler creates a new scenario handler.
func newScenarioHandler(scenarioService service.Scenario, opts ...Opts) *ScenarioHandler {
	options := defaultOptions()
	for _, o := range opts {
		o(options)
	}
	return &ScenarioHandler{
		scenarioService: scenarioService,
		logger:          options.Logger,
	}
}

// CreateScenario create a scenario for a project.
func (h *ScenarioHandler) CreateScenario(ctx echo.Context, projectID uuid.UUID) error {
	httpScenario := &api.CreateScenarioJSONRequestBody{}
	err := json.NewDecoder(ctx.Request().Body).Decode(&httpScenario)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid create scenario payload")
	}
	createScenarioRequest := &app.CreateScenarioRequest{
		Name:      httpScenario.Name,
		SpecType:  app.ScenarioSpecType(httpScenario.SpecType),
		Spec:      httpScenario.Spec,
		ProjectID: projectID,
	}
	scenario, err := h.scenarioService.CreateScenario(ctx.Request().Context(), createScenarioRequest)
	switch {
	case errors.Is(err, app.ErrScenarioAlreadyExists):
		return echo.NewHTTPError(http.StatusConflict, fmt.Sprintf("scenario with name %s already exists for given project", httpScenario.Name))
	case errors.Is(err, app.ErrProjectNotFound):
		return echo.NewHTTPError(http.StatusNotFound, "project not found")
	case err != nil:
		h.logger.Error("unable to create scenario", slog.String("error", err.Error()))
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to create scenario")
	}

	return ctx.JSON(http.StatusCreated, appScenarioToHTTPScenario(scenario))
}

// ListScenariosForProject lists all scenarios for a project.
func (h *ScenarioHandler) ListScenariosForProject(ctx echo.Context, projectId uuid.UUID, params api.ListScenariosForProjectParams) error {
	listScenariosRequest := &app.ListScenariosRequest{
		Limit:     100,
		Offset:    0,
		ProjectID: projectId,
	}
	if params.Limit != nil {
		listScenariosRequest.Limit = *params.Limit
	}
	if params.Offset != nil {
		listScenariosRequest.Offset = *params.Offset
	}

	scenarios, err := h.scenarioService.ListScenarios(ctx.Request().Context(), listScenariosRequest)
	if err != nil {
		h.logger.Error("unable to get scenarios for project", slog.String("error", err.Error()))
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to get scenarios for project")
	}

	result := make([]api.Scenario, len(scenarios))
	for i, scenario := range scenarios {
		result[i] = appScenarioToHTTPScenario(scenario)
	}

	return ctx.JSON(http.StatusOK, result)
}

func appScenarioToHTTPScenario(scenario *app.Scenario) api.Scenario {
	return api.Scenario{
		ID:        scenario.ID,
		Name:      scenario.Name,
		Spec:      scenario.Spec,
		SpecType:  api.ScenarioSpecType(scenario.SpecType.String()),
		ProjectID: scenario.ProjectID,
	}
}
