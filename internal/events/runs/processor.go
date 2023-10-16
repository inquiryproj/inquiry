package runs

import (
	"github.com/google/uuid"

	"github.com/inquiryproj/inquiry/internal/app"
)

// Processor processes runs.
type Processor interface {
	// Process processes a run.
	Process(projectID uuid.UUID) (*app.ProjectRunOutput, error)
}
