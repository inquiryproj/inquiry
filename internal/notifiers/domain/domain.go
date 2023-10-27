// Package domain holds the notifier domain objects.
package domain

import "time"

// ProjectRun is the output of a project run.
type ProjectRun struct {
	Name         string
	Success      bool
	Version      string
	Duration     time.Duration
	Time         time.Time
	ScenarioRuns []*ScenarioRunDetails
}

// ScenarioRunDetails is the output of a scenario run.
type ScenarioRunDetails struct {
	Name                 string
	Duration             time.Duration
	SuccessfulAssertions int
	Assertions           int
	Success              bool
}
