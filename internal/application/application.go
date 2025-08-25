// Package application provide the application wrapper.
package application

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hoppermq/hopper/internal/config"
	"github.com/hoppermq/hopper/pkg/domain"
)

// Application is the application structure wrapper.
type Application struct {
	configuration config.Configuration
	logger        *slog.Logger

	service domain.Service
	eb      domain.IEventBus
	running chan bool
	stop    chan os.Signal
}

// Option is the function that configure the service.
type Option func(*Application)

// WithLogger set the logger to the application.
func WithLogger(logger *slog.Logger) Option {
	return func(a *Application) {
		a.logger = logger
	}
}

func WithService(service domain.Service) Option {
	return func(a *Application) {
		a.service = service
	}
}

func WithEventBus(eb domain.IEventBus) Option {
	return func(a *Application) {
		a.eb = eb
	}
}

func WithConfiguration(cfg *config.Configuration) Option {
	return func(a *Application) {
		a.configuration = *cfg
	}
}

// New create a new application instance.
func New(opts ...Option) *Application {
	app := &Application{
		running: make(chan bool, 1),
		stop:    make(chan os.Signal, 1),
	}

	for _, opt := range opts {
		opt(app)
	}

	return app
}

func (a *Application) Start() {
	a.logger.Info(
		"Application STARTED",
		"name", a.configuration.App.Name,
		"version", a.configuration.App.Version,
		"id", a.configuration.App.ID,
	)
	ctx := context.Background()
	signal.Notify(a.stop, syscall.SIGINT, syscall.SIGTERM)

	if eventBusAware, ok := a.service.(domain.EventBusAware); ok {
		eventBusAware.RegisterEventBus(a.eb)
	}

	go func() {
		if err := a.service.Run(ctx); err != nil {
			a.logger.Error("Failed to start component: ", a.service.Name(), err)
			a.stop <- syscall.SIGTERM
		}
	}()

	a.running <- true
	<-a.stop

	a.logger.Info("Shutting down application")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := a.service.Stop(shutdownCtx); err != nil {
		a.logger.Error("Failed to stop service: ", a.service.Name(), err)
	}
	a.logger.Info("Application STOPPED")
}
