package handlers

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/inquiryproj/inquiry/internal/service"
)

// RunHandler handles run requests.
type RunHandler struct {
	runnerService service.Runner
	logger        *slog.Logger
}

// newRunHandler creates a new run handler.
func newRunHandler(runnerService service.Runner, options *Options) *RunHandler {
	return &RunHandler{
		runnerService: runnerService,
		logger:        options.Logger,
	}
}

// RunProject runs all scenarios for a given project.
func (h *RunHandler) RunProject(ctx echo.Context, id uuid.UUID) error {
	_ = id

	_ = ctx
	return nil
}
