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
	configuration *config.Configuration
	logger        *slog.Logger

	services []domain.Service
	eb       domain.IEventBus
	running  chan bool
	stop     chan os.Signal
}

// Option is the function that configure the service.
type Option func(*Application)

// WithLogger set the logger to the application.
func WithLogger(logger *slog.Logger) Option {
	return func(a *Application) {
		a.logger = logger
	}
}

func WithService(services ...domain.Service) Option {
	return func(a *Application) {
		a.services = append(a.services, services...)
	}
}

func WithEventBus(eb domain.IEventBus) Option {
	return func(a *Application) {
		a.eb = eb
	}
}

// WithConfiguration inject the configuration to the application.
func WithConfiguration(cfg *config.Configuration) Option {
	return func(a *Application) {
		a.configuration = cfg
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

	for _, s := range a.services {
		if eventBusAware, ok := s.(domain.EventBusAware); ok {
			eventBusAware.RegisterEventBus(a.eb)
		}

		go func(svc domain.Service) {
			if err := s.Run(ctx); err != nil {
				a.logger.Error("Failed to start component: ", s.Name(), err)
				a.stop <- syscall.SIGTERM
			}
		}(s)

	}

	a.running <- true
	<-a.stop

	a.logger.Info("Shutting down application")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, s := range a.services {
		if err := s.Stop(shutdownCtx); err != nil {
			a.logger.Error("Failed to stop service: ", s.Name(), err)
		}
	}
	a.logger.Info("Application STOPPED")
}
