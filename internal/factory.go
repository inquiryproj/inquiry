// Package factory instantiates the application.
package factory

import (
	"log/slog"
	"os"

	server "github.com/inquiryproj/inquiry/internal/http"
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

	serviceWrapper := serviceFactory(repositoryWrapper)

	handlerWrapper := handlers.NewHandlerWrapper(serviceWrapper,
		handlers.WithLogger(logger),
	)

	return server.NewAPI(handlerWrapper,
		server.WithLogger(logger),
		server.WithPort(cfg.ServerConfig.Port),
		server.WithShutdownDelay(cfg.ServerConfig.ShutdownDelay),
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

func serviceFactory(repositoryWrapper repository.Wrapper) service.Wrapper {
	return service.NewServiceWrapper(repositoryWrapper)
}

func repositoryFactory(repositoryConfig RepositoryConfig) (repository.Wrapper, error) {
	return repository.NewWrapper(
		repository.WithType(repositoryConfig.RepositoryType.String()),
		repository.WithDSN(repositoryConfig.DSN),
	)
}
