package completions

import (
	"context"

	"github.com/google/uuid"

	"github.com/inquiryproj/inquiry/internal/app"
	"github.com/inquiryproj/inquiry/internal/integrations/slack"
	"github.com/inquiryproj/inquiry/internal/notifiers"
)

// Integration is an integration which can send completions.
type Integration interface {
	// SendCompletion sends a completion message.
	SendCompletion(ctx context.Context, projectRun notifiers.ProjectRun) error
}

// Processor processes runs.
type Processor interface {
	// Process processes a run.
	Process(runID uuid.UUID) (*app.ProjectRunOutput, error)
}

type processor struct {
	slackClient slack.Client
	test        []Integration
}
