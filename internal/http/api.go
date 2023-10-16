package http

import (
	"context"
	"errors"
	"fmt"
	"log"
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

// Options represents the options for the API.
type Options struct {
	Port          int
	ShutdownDelay time.Duration
	Logger        *slog.Logger
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

// API is the API server.
type API struct {
	e *echo.Echo

	port          int
	shutdownDelay time.Duration
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
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	return &API{
		e:             e,
		port:          options.Port,
		shutdownDelay: options.ShutdownDelay,
	}
}

// Run runs the API server and handles graceful shutdown.
func (a *API) Run() error {
	errChan := make(chan error, 1)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		// Start server

		err := a.e.Start(fmt.Sprintf(":%d", a.port))
		if errors.Is(http.ErrServerClosed, err) {
			log.Default().Println("server closed gracefully")
		} else if err != nil {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		return err
	case <-c:
	}
	time.Sleep(a.shutdownDelay)
	log.Default().Println("Shutting down server")
	err := a.e.Shutdown(context.Background())
	if err != nil {
		log.Default().Println("unable to shutdown server")
	}
	return nil
}
