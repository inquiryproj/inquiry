package http

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	httpMocks "github.com/inquiryproj/inquiry/internal/http/mocks"
)

func TestNewAPIDefault(t *testing.T) {
	serverImpl := httpMocks.NewServerInterface(t)
	api := NewAPI(serverImpl)
	assert.Equal(t, 3000, api.port)
	assert.Equal(t, time.Duration(0), api.shutdownDelay)
	assert.Len(t, api.runnables, 0)
	assert.NotNil(t, api.logger)
	assert.NotNil(t, api.e)
	assert.NotNil(t, api.errChan)
	assert.NotNil(t, api.shutDownChan)
}

func TestNewAPIWithOpts(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	serverImpl := httpMocks.NewServerInterface(t)
	runnableMock := httpMocks.NewRunnable(t)
	api := NewAPI(serverImpl,
		WithPort(42),
		WithShutdownDelay(1*time.Second),
		WithLogger(logger),
		WithRunnable(runnableMock),
	)
	assert.Equal(t, 42, api.port)
	assert.Equal(t, time.Second*1, api.shutdownDelay)
	assert.Len(t, api.runnables, 1)
	assert.Equal(t, logger, api.logger)
	assert.NotNil(t, api.e)
	assert.NotNil(t, api.errChan)
}

func TestStartRunnables(t *testing.T) {
	serverImpl := httpMocks.NewServerInterface(t)
	runnableMock := httpMocks.NewRunnable(t)
	done := make(chan struct{})
	runnableMock.On("Start").Return(nil).Run(func(args mock.Arguments) {
		done <- struct{}{}
	}).Times(1)
	runnableMock.On("Name").Return("mock").Times(1)
	api := NewAPI(serverImpl,
		WithRunnable(runnableMock),
	)
	api.startRunnables()
	<-done
}

func TestUnableToStartRunnables(t *testing.T) {
	serverImpl := httpMocks.NewServerInterface(t)
	runnableMock := httpMocks.NewRunnable(t)
	runnableMock.On("Start").Return(assert.AnError).Times(1)
	runnableMock.On("Name").Return("mock").Times(2)
	api := NewAPI(serverImpl,
		WithRunnable(runnableMock),
	)
	api.startRunnables()
	err := <-api.errChan
	assert.Error(t, err)
}

func TestShutDownServer(t *testing.T) {
	runnableMock := httpMocks.NewRunnable(t)
	runnableMock.On("Shutdown", mock.Anything).Return(nil).Times(1)
	runnableMock.On("Name").Return("mock").Times(1)

	serverImpl := httpMocks.NewServerInterface(t)
	api := NewAPI(serverImpl,
		WithRunnable(runnableMock),
	)
	err := api.shutDownServer()
	assert.NoError(t, err)
}

func TestRunAndShutdown(t *testing.T) {
	runnableMock := httpMocks.NewRunnable(t)

	runnableMock.On("Shutdown", mock.Anything).Return(nil).Times(1)
	runnableMock.On("Name").Return("mock").Times(2)

	serverImpl := httpMocks.NewServerInterface(t)
	api := NewAPI(serverImpl,
		WithRunnable(runnableMock),
	)
	runnableMock.On("Start").Return(nil).Run(func(args mock.Arguments) {
		api.shutDownChan <- os.Kill
	}).Times(1)
	err := api.Run()
	assert.NoError(t, err)
}

func TestOptionsNoAuth(t *testing.T) {
	options := defaultOptions()
	for _, opt := range []Opts{
		WithPort(42),
		WithShutdownDelay(1 * time.Second),
		WithLogger(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))),
		WithRunnable(httpMocks.NewRunnable(t)),
		WithAuthDisabled(),
	} {
		opt(options)
	}
	assert.Equal(t, 42, options.Port)
	assert.Equal(t, time.Second*1, options.ShutdownDelay)
	assert.Len(t, options.Runnables, 1)
	assert.NotNil(t, options.Logger)
	assert.NotNil(t, options.APIKeyRepository)
	assert.False(t, options.WithAuthEnabled)
}

func TestOptionsAuth(t *testing.T) {
	options := defaultOptions()
	for _, opt := range []Opts{
		WithPort(42),
		WithShutdownDelay(1 * time.Second),
		WithLogger(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))),
		WithRunnable(httpMocks.NewRunnable(t)),
		WithAuthEnabled(&apiKeyDenyAll{}),
	} {
		opt(options)
	}
	assert.Equal(t, 42, options.Port)
	assert.Equal(t, time.Second*1, options.ShutdownDelay)
	assert.Len(t, options.Runnables, 1)
	assert.NotNil(t, options.Logger)
	assert.NotNil(t, options.APIKeyRepository)
	assert.True(t, options.WithAuthEnabled)
}
