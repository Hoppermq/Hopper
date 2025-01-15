// Package application provide the application wrapper.
package application

import (
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
  a.logger.Info("HOPPER STARTED")
  signal.Notify(a.stop, syscall.SIGINT, syscall.SIGTERM)
  a.running <- true;
  <- a.stop;
}
