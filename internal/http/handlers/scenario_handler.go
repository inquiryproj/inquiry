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
	httpInternal "github.com/inquiryproj/inquiry/internal/http"
	"github.com/inquiryproj/inquiry/internal/service"
)

// ScenarioHandler handles scenario requests.
type ScenarioHandler struct {
	scenarioService service.Scenario
	logger          *slog.Logger
}

// newScenarioHandler creates a new scenario handler.
func newScenarioHandler(scenarioService service.Scenario, options *Options) *ScenarioHandler {
	return &ScenarioHandler{
		scenarioService: scenarioService,
		logger:          options.Logger,
	}
}

// CreateScenario create a scenario for a project.
func (h *ScenarioHandler) CreateScenario(ctx echo.Context, id uuid.UUID) error {
	httpScenario := &httpInternal.Scenario{}
	err := json.NewDecoder(ctx.Request().Body).Decode(&httpScenario)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid create scenario payload")
	}
	createScenarioRequest := &app.CreateScenarioRequest{
		Name:      httpScenario.Name,
		SpecType:  app.ScenarioSpecType(httpScenario.SpecType),
		Spec:      httpScenario.Spec,
		ProjectID: id,
	}
	scenario, err := h.scenarioService.CreateScenario(ctx.Request().Context(), createScenarioRequest)
	if errors.Is(err, app.ErrScenarioAlreadyExists) {
		return echo.NewHTTPError(http.StatusConflict, fmt.Sprintf("scenario with name %s already exists for given project", httpScenario.Name))
	} else if err != nil {
		h.logger.Error("unable to create scenario", slog.String("error", err.Error()))
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to create scenario")
	}

	return ctx.JSON(http.StatusCreated, httpInternal.Scenario{
		ID:        scenario.ID,
		Name:      scenario.Name,
		Spec:      scenario.Spec,
		SpecType:  scenario.SpecType.String(),
		ProjectID: scenario.ProjectID,
	})
}
