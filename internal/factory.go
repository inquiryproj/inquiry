// Package factory instantiates the application.
package factory

import (
	"log/slog"
	"os"

	"github.com/google/uuid"

	"github.com/inquiryproj/inquiry/internal/events"
	"github.com/inquiryproj/inquiry/internal/events/completions"
	"github.com/inquiryproj/inquiry/internal/events/runs"
	"github.com/inquiryproj/inquiry/internal/http"
	"github.com/inquiryproj/inquiry/internal/http/handlers"
	"github.com/inquiryproj/inquiry/internal/notifiers"
	"github.com/inquiryproj/inquiry/internal/repository"
	"github.com/inquiryproj/inquiry/internal/service"
)

// App is the application.
type App interface {
	Run() error
}

// NewApp creates a new App instance.
func NewApp() (App, error) {
	cfg, err := NewConfig()
	if err != nil {
		return nil, err
	}
	logger := loggerFactory(cfg.LogLevel, cfg.LogFormat)

	notifierServices := notifiersFactory(cfg.NotifiersConfig)

	repositoryWrapper, err := repositoryFactory(cfg.RepositoryConfig, cfg.ServerConfig.APIKey)
	if err != nil {
		logger.Error("failed to initialise repository", slog.String("error", err.Error()))
		return nil, err
	}

	completionsProducer, completionsConsumer, err := completionEventsFactory(notifierServices, repositoryWrapper)
	if err != nil {
		logger.Error("failed to initialise runs events", slog.String("error", err.Error()))
		return nil, err
	}

	runsProducer, runsConsumer, err := runEventsFactory(completionsProducer, repositoryWrapper)
	if err != nil {
		logger.Error("failed to initialise runs events", slog.String("error", err.Error()))
		return nil, err
	}

	serviceWrapper := serviceFactory(repositoryWrapper, runsProducer)

	handlerWrapper := handlers.NewHandlerWrapper(serviceWrapper,
		handlers.WithLogger(logger),
	)

	opts := []http.Opts{
		http.WithLogger(logger),
		http.WithPort(cfg.ServerConfig.Port),
		http.WithShutdownDelay(cfg.ServerConfig.ShutdownDelay),
		http.WithRunnable(runsConsumer),
		http.WithRunnable(completionsConsumer),
	}
	if cfg.ServerConfig.AuthEnabled {
		opts = append(opts, http.WithAuthEnabled(repositoryWrapper.APIKey))
	} else {
		opts = append(opts, http.WithAuthDisabled())
	}

	return http.NewAPI(handlerWrapper,
		opts...,
	), nil
}

func loggerFactory(logLevel LogLevel, logFormat LogFormat) *slog.Logger {
	switch logFormat {
	case LogFormatJSON:
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: leveler(logLevel),
		}))
	case LogFormatText:
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: leveler(logLevel),
		}))
	default:
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	}
}

func leveler(logLevel LogLevel) slog.Leveler {
	switch logLevel {
	case LogLevelDebug:
		return slog.LevelDebug
	case LogLevelInfo:
		return slog.LevelInfo
	case LogLevelWarn:
		return slog.LevelWarn
	case LogLevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func notifiersFactory(notifiersConfig NotifiersConfig) []notifiers.Notifier {
	notifierOpts := []notifiers.Opts{}
	if notifiersConfig.SlackConfig.WebhookURL != "" {
		notifierOpts = append(notifierOpts, notifiers.WithSlackEnabled(notifiersConfig.SlackConfig.WebhookURL))
	}
	return notifiers.NewNotifiers(notifierOpts...)
}

func completionEventsFactory(notifierServices []notifiers.Notifier, repositoryWrapper *repository.Wrapper) (events.Producer[uuid.UUID], http.Runnable, error) {
	completionProcessor := completionProcessorFactory(notifierServices, repositoryWrapper)
	producer, consumer, err := completions.NewProducerConsumer(completionProcessor)
	if err != nil {
		return nil, nil, err
	}
	return producer, newRunnableConsumer(consumer, "completion consumer"), nil
}

func runEventsFactory(completionsProducer events.Producer[uuid.UUID], repositoryWrapper *repository.Wrapper) (events.Producer[uuid.UUID], http.Runnable, error) {
	runProcessor := runProcessorFactory(completionsProducer, repositoryWrapper)
	producer, consumer, err := runs.NewProducerConsumer(runProcessor)
	if err != nil {
		return nil, nil, err
	}
	return producer, newRunnableConsumer(consumer, "runs consumer"), nil
}

type runnableConsumer struct {
	events.Consumer
	name string
}

func newRunnableConsumer(consumer events.Consumer, name string) http.Runnable {
	return &runnableConsumer{
		Consumer: consumer,
		name:     name,
	}
}

func (r *runnableConsumer) Start() error {
	return r.Consume()
}

func (r *runnableConsumer) Name() string {
	return r.name
}

func completionProcessorFactory(notifierServices []notifiers.Notifier, repositoryWrapper *repository.Wrapper) runs.Processor {
	return completions.NewProcessor(notifierServices, repositoryWrapper.Run, repositoryWrapper.Project)
}

func runProcessorFactory(completionsProducer events.Producer[uuid.UUID], repositoryWrapper *repository.Wrapper) runs.Processor {
	return runs.NewProcessor(completionsProducer, repositoryWrapper.Scenario, repositoryWrapper.Run)
}

func serviceFactory(repositoryWrapper *repository.Wrapper, runsProducer events.Producer[uuid.UUID]) service.Wrapper {
	return service.NewServiceWrapper(repositoryWrapper, runsProducer)
}

func repositoryFactory(repositoryConfig RepositoryConfig, apiKey string) (*repository.Wrapper, error) {
	return repository.NewWrapper(
		repository.WithType(repositoryConfig.RepositoryType.String()),
		repository.WithDSN(repositoryConfig.DSN),
		repository.WithAPIKey(apiKey),
	)
}
