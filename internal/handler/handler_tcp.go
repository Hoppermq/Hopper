package handler

import (
	"context"
	"log/slog"
	"net"
)

// TCP is an TCP handler
type TCP struct {
  Listener net.Listener
  logger *slog.Logger
}


type config struct {
  lconf net.ListenConfig
  ctx context.Context
  logger *slog.Logger
}

type Option func(*config) error

func WithLogger(logger *slog.Logger) Option {
  return func(c *config) error {
    c.logger = logger

    return nil
  }
}

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
    logger: handlerConfig.logger,
  }, nil
}

func (t *TCP) handleConnection() error {
  for {
    conn, err := t.Listener.Accept();
    if err != nil {
      t.logger.Error("failed to accept connection", err)
      return err
    }
    go t.processConnection(conn);
  }

}

func (t *TCP) processConnection(conn net.Conn) {
  defer conn.Close()
}

func (t *TCP) Start(_ context.Context) error {
  go t.handleConnection();
  return nil
}

func (t *TCP) Stop() error {
  return t.Listener.Close() 
}
