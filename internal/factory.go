// Package factory instantiates the application.
package factory

import (
	"log/slog"
	"os"

	"github.com/inquiryproj/inquiry/internal/events/runs"
	"github.com/inquiryproj/inquiry/internal/events/runs/run"
	"github.com/inquiryproj/inquiry/internal/http"
	"github.com/inquiryproj/inquiry/internal/http/handlers"
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

	repositoryWrapper, err := repositoryFactory(cfg.RepositoryConfig)
	if err != nil {
		logger.Error("failed to initialise repository", slog.String("error", err.Error()))
		return nil, err
	}

	runsProducer, runsConsumer, err := runEventsFactory(repositoryWrapper)
	if err != nil {
		logger.Error("failed to initialise runs events", slog.String("error", err.Error()))
		return nil, err
	}

	serviceWrapper := serviceFactory(repositoryWrapper, runsProducer)

	handlerWrapper := handlers.NewHandlerWrapper(serviceWrapper,
		handlers.WithLogger(logger),
	)

	return http.NewAPI(handlerWrapper,
		http.WithLogger(logger),
		http.WithPort(cfg.ServerConfig.Port),
		http.WithShutdownDelay(cfg.ServerConfig.ShutdownDelay),
		http.WithRunnable(runsConsumer),
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

func runEventsFactory(repositoryWrapper *repository.Wrapper) (runs.Producer, http.Runnable, error) {
	runProcessor := processorFactory(repositoryWrapper)
	producer, consumer, err := runs.NewProducerConsumer(runProcessor)
	if err != nil {
		return nil, nil, err
	}
	return producer, newRunnableConsumer(consumer, "runs consumer"), nil
}

type runnableConsumer struct {
	runs.Consumer
	name string
}

func newRunnableConsumer(consumer runs.Consumer, name string) http.Runnable {
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

func processorFactory(repositoryWrapper *repository.Wrapper) run.Processor {
	return run.NewProcessor(repositoryWrapper.Scenario)
}

func serviceFactory(repositoryWrapper *repository.Wrapper, runsProducer runs.Producer) service.Wrapper {
	return service.NewServiceWrapper(repositoryWrapper, runsProducer)
}

func repositoryFactory(repositoryConfig RepositoryConfig) (*repository.Wrapper, error) {
	return repository.NewWrapper(
		repository.WithType(repositoryConfig.RepositoryType.String()),
		repository.WithDSN(repositoryConfig.DSN),
	)
}
