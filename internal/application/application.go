// Package application provide the application wrapper.
package application

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

// Application is the application structure wrapper.
type Application struct {
	name    string
	version string
	id      string

	logger *slog.Logger

  service Service
  running chan bool
  stop chan os.Signal
}

// Option is the function that configure the service.
type Option func(*Application);

// WithLogger set the logger to the application.
func WithLogger(logger *slog.Logger) Option {
  return func(a *Application) {
    a.logger = logger
  }
}

func WithService(service Service) Option {
  return func(a *Application) {
    a.service = service
  }
}

// New create a new application instance.
func New(opts ...Option) *Application {
  app := &Application{
    name: "Hopper",
    version: "v/0.0.1",
    id: "hppr-id-01",
    running: make(chan bool, 1),
    stop: make(chan os.Signal, 1),
  }

  for _, opt := range opts {
    opt(app)
  }
  
  return app 
}

func (a *Application) Start() {
  a.logger.Info("Application STARTED")
  ctx := context.Background();
  signal.Notify(a.stop, syscall.SIGINT, syscall.SIGTERM)

  go func() {
    if err := a.service.Run(ctx); err != nil {
      a.logger.Error("Failed to start component: ", a.service.Name(), err)
      a.stop <- syscall.SIGTERM
    }
  }()

  a.running <- true;
  <- a.stop;
}
