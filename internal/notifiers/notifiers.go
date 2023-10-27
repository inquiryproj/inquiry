// Package notifiers provides domain models for different notifier integrations.
package notifiers

import (
	"context"

	"github.com/inquiryproj/inquiry/internal/notifiers/domain"
	"github.com/inquiryproj/inquiry/internal/notifiers/slack"
)

// Notifier is a notifier which can send completions.
type Notifier interface {
	// SendCompletion sends a completion message.
	SendCompletion(ctx context.Context, projectRun *domain.ProjectRun) error
}

type options struct {
	SlackEnabled    bool
	SlackWebhookURL string
}

func defaultOptions() *options {
	return &options{
		SlackEnabled:    false,
		SlackWebhookURL: "",
	}
}

// Opts represents a function that modifies the options.
type Opts func(*options)

// WithSlackEnabled enables slack.
func WithSlackEnabled(webhookURL string) Opts {
	return func(o *options) {
		o.SlackEnabled = true
		o.SlackWebhookURL = webhookURL
	}
}

// NewNotifiers creates a new set of notifiers.
func NewNotifiers(opts ...Opts) []Notifier {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	notifiers := []Notifier{}
	if options.SlackEnabled {
		notifiers = append(notifiers, slack.NewClient(options.SlackWebhookURL))
	}
	return notifiers
}
