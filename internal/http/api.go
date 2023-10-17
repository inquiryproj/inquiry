package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
)

// Runnable is the interface for runnable components managed by the API.
type Runnable interface {
	Start() error
	Name() string
	Shutdown(ctx context.Context) error
}

// Options represents the options for the API.
type Options struct {
	Port          int
	ShutdownDelay time.Duration
	Logger        *slog.Logger
	Runnables     []Runnable
}

func defaultOptions() *Options {
	return &Options{
		Port:          3000,
		ShutdownDelay: time.Duration(0),
		Logger:        slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})),
	}
}

// Opts represents a function that modifies the options.
type Opts func(*Options)

// WithPort sets the port.
func WithPort(port int) Opts {
	return func(o *Options) {
		o.Port = port
	}
}

// WithShutdownDelay sets the shutdown delay.
func WithShutdownDelay(delay time.Duration) Opts {
	return func(o *Options) {
		o.ShutdownDelay = delay
	}
}

// WithLogger sets the logger.
func WithLogger(logger *slog.Logger) Opts {
	return func(o *Options) {
		o.Logger = logger
	}
}

// WithRunnable adds a runnable.
func WithRunnable(runnable Runnable) Opts {
	return func(o *Options) {
		o.Runnables = append(o.Runnables, runnable)
	}
}

// API is the API server.
type API struct {
	e *echo.Echo

	port          int
	shutdownDelay time.Duration
	runnables     []Runnable

	errChan chan error

	logger *slog.Logger
}

// NewAPI creates a new API server.
func NewAPI(handler ServerInterface, opts ...Opts) *API {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(
		slogecho.NewWithConfig(
			options.Logger,
			slogecho.Config{
				DefaultLevel:     slog.LevelInfo,
				ClientErrorLevel: slog.LevelWarn,
				ServerErrorLevel: slog.LevelError,
				WithRequestID:    false,
				Filters: []slogecho.Filter{
					slogecho.IgnoreStatus(http.StatusOK, http.StatusNotFound, http.StatusUnauthorized),
				},
			},
		),
	)
	RegisterHandlers(e, handler)

	// Middlewares
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	return &API{
		e:             e,
		port:          options.Port,
		shutdownDelay: options.ShutdownDelay,
		runnables:     options.Runnables,
		logger:        options.Logger,

		errChan: make(chan error, 1),
	}
}

// Run runs the API server and handles graceful shutdown.
func (a *API) Run() error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go a.startHTTPServer()
	go a.startRunnables()

	select {
	case err := <-a.errChan:
		return err
	case <-c:
	}

	return a.shutDownServer()
}

func (a *API) startHTTPServer() {
	a.logger.Info("starting server")
	err := a.e.Start(fmt.Sprintf(":%d", a.port))
	if errors.Is(http.ErrServerClosed, err) {
		a.logger.Info("server closed gracefully")
	} else if err != nil {
		a.errChan <- err
	}
}

func (a *API) startRunnables() {
	for _, runnable := range a.runnables {
		go func(runnable Runnable) {
			a.logger.Info("starting component", slog.String("runnable_name", runnable.Name()))
			err := runnable.Start()
			if err != nil {
				a.logger.Error("unable to start component", slog.String("runnable_name", runnable.Name()), slog.String("error", err.Error()))
				a.errChan <- err
			}
		}(runnable)
	}
}

func (a *API) shutDownServer() error {
	for _, runnable := range a.runnables {
		a.logger.Info("shutting down component", slog.String("runnable_name", runnable.Name()))
		err := runnable.Shutdown(context.Background())
		if err != nil {
			a.logger.Error("unable to shutdown component", slog.String("runnable_name", runnable.Name()), slog.String("error", err.Error()))
		}
	}
	time.Sleep(a.shutdownDelay)
	a.logger.Info("Shutting down server")
	err := a.e.Shutdown(context.Background())
	if err != nil {
		a.logger.Error("unable to shutdown server", slog.String("error", err.Error()))
		return err
	}
	return nil
}
