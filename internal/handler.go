package internal

import (
	"context"
	"log/slog"
	"net"
)

// TCP is an TCP handler
type TCP struct {
  Listener net.Listener
  logger slog.Logger
}


type config struct {
  lconf net.ListenConfig
  ctx context.Context
}

type Option func(*config) error

func WithListener(listenerConfig net.ListenConfig) Option {
  return func(c *config) error {
    c.lconf = listenerConfig;

    return nil
  }
}

func WithContext(ctx context.Context) Option {
  return func(c *config) error {
    c.ctx = ctx;

    return nil
  }
}

// NewTCP return the new tcp handler
func NewTCP(opts ...Option) (*TCP, error) {
  handlerConfig := &config{}
  for _, opt := range opts {
    opt(handlerConfig)
  }

  l, err := handlerConfig.lconf.Listen(handlerConfig.ctx, "tcp", ":9091")
  if err != nil {
    return nil, err
  }

  return &TCP{
    Listener: l,
  }, nil
}

// here we will handle all command injected in the register
