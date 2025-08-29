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
	"github.com/hoppermq/hopper/internal/events"
	"github.com/hoppermq/hopper/pkg/domain"
)

// Application is the application structure wrapper.
type Application struct {
	configuration *config.Configuration
	logger        *slog.Logger

	services []domain.IService
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

func WithService(services ...domain.IService) Option {
	return func(a *Application) {
		a.services = append(a.services, services...)
	}
}

func WithEventBus(opts ...events.Option) Option {
	return func(a *Application) {
		a.eb = events.NewEventBus(1000) //should take the event opts config.
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
		go func(svc domain.IService) {
			svc.RegisterEventBus(a.eb)

			if err := svc.Run(ctx); err != nil {
				a.logger.Error("Failed to start component: ", s.Name(), err)
				a.stop <- syscall.SIGTERM
			}
		}(s)

	}

	a.running <- true
	a.Stop()
	a.logger.Info("application shutted down succesfully")
}

func (a *Application) Stop() {
	a.logger.Info(
		"shutting down application",
		"name",
		a.getName(),
		"application_id",
		a.getID(),
	)
	<-a.stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, s := range a.services {
		if err := s.Stop(shutdownCtx); err != nil {
			a.logger.Error("failed to stop service: ", s.Name(), err)
		}
	}
}

func (a *Application) getName() string {
	return a.configuration.App.Name
}

func (a *Application) getID() domain.ID {
	return domain.ID(a.configuration.App.ID)
}
