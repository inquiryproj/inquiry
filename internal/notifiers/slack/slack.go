// Package slack provides a slack client for reporting test results.
package slack

import (
	"context"
	"fmt"
	"time"

	"github.com/slack-go/slack"

	"github.com/inquiryproj/inquiry/internal/notifiers/domain"
)

const (
	slackUserName = "Inquiry"
)

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

// Client is the interface for sending slack messages.
type Client interface {
	SendCompletion(ctx context.Context, projectRun *domain.ProjectRun) error
}

// NewClient creates a new slack client.
func NewClient(webhookURL string) Client {
	return &client{
		webhookURL: webhookURL,
	}
}

type client struct {
	webhookURL string
}

// SendProjectRun sends a project run to slack.
func (c *client) SendCompletion(ctx context.Context, projectRun *domain.ProjectRun) error {
	// FIXME we might want to be able to send to different webhooks for different projects.
	return slack.PostWebhookContext(ctx, c.webhookURL, &slack.WebhookMessage{
		Username: slackUserName,
		IconURL:  "https://github.com/inquiryproj/inquiry/blob/main/assets/logo.png?raw=true",
		Blocks:   buildSlackMessageBlock(projectRun),
	},
	)
}

func buildSlackMessageBlock(projectRun *domain.ProjectRun) *slack.Blocks {
	return &slack.Blocks{
		BlockSet: []slack.Block{
			buildHeader(projectRun),
			buildRunOverviewSection(projectRun),
			buildScenarioDetailsSection(projectRun),
		},
	}
}

func buildHeader(projectRun *domain.ProjectRun) *slack.HeaderBlock {
	if projectRun.Success {
		return slack.NewHeaderBlock(slack.NewTextBlockObject("plain_text", fmt.Sprintf("Project %s success  :white_check_mark:", projectRun.Name), false, false))
	}
	return slack.NewHeaderBlock(slack.NewTextBlockObject("plain_text", fmt.Sprintf("Project %s failed :x:", projectRun.Name), false, false))
}

func buildRunOverviewSection(projectRun *domain.ProjectRun) *slack.SectionBlock {
	return slack.NewSectionBlock(
		nil,
		[]*slack.TextBlockObject{
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Version:*\n%s", projectRun.Version), false, false),
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Duration:*\n%s", projectRun.Duration.String()), false, false),
		},
		nil,
	)
}

func buildScenarioDetailsSection(projectRun *domain.ProjectRun) *slack.SectionBlock {
	scenarios := ""
	assertions := ""
	for _, s := range projectRun.ScenarioRuns {
		scenarioText := fmt.Sprintf(":x: %s\n", s.Name)
		if s.Success {
			scenarioText = fmt.Sprintf(":white_check_mark: %s\n", s.Name)
		}
		scenarios += fmt.Sprintf("%s\n", scenarioText)
		assertions += fmt.Sprintf("%d/%d\n", s.SuccessfulAssertions, s.Assertions)
	}
	return slack.NewSectionBlock(
		nil,
		[]*slack.TextBlockObject{
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Scenarios:*\n%s", scenarios), false, false),
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Assertions:*\n%s", assertions), false, false),
		},
		nil,
	)
}
