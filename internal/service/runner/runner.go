// Package runner implements the runner service.
package runner

import (
	"context"

	"github.com/inquiryproj/inquiry/internal/app"
)

// Runner is the runner service.
type Runner struct{}

// NewService initialises the runner service.
func NewService() *Runner {
	return &Runner{}
}

// RunProject runs all scenarios for a given project.
func (s *Runner) RunProject(ctx context.Context, runProjectRequest *app.RunProjectRequest) (*app.ProjectRunOutput, error) {
	// FIXME: implement
	_ = ctx
	_ = runProjectRequest
	return nil, nil
}
