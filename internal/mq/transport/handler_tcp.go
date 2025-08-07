// Package transport provides all the transport handler for the Hopper message queue system.
package transport

import (
	"context"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/hoppermq/hopper/internal/mq/core"
)

// TCP is an TCP handler
type TCP struct {
	Listener net.Listener
	logger   *slog.Logger

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	broker *core.Broker
	cm     *core.ClientManager
}

type config struct {
	lconf  net.ListenConfig
	ctx    context.Context
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
		c.lconf = listenerConfig

		return nil
	}
}

func WithContext(ctx context.Context) Option {
	return func(c *config) error {
		c.ctx = ctx

		return nil
	}
}

// NewTCP return the new tcp handler
func NewTCP(opts ...Option) (*TCP, error) {
	handlerConfig := &config{}
	for _, opt := range opts {
		err := opt(handlerConfig)
		if err != nil {
			return nil, err
		}
	}

	l, err := handlerConfig.lconf.Listen(handlerConfig.ctx, "tcp", ":9091")
	if err != nil {
		return nil, err
	}

	return &TCP{
		Listener: l,
		logger:   handlerConfig.logger,
	}, nil
}

func (t *TCP) HandleConnection(ctx context.Context) error {
	for {
		conn, err := t.Listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				t.logger.Info("context cancelled, stopping connection handler")
				return ctx.Err()
			default:
				t.logger.Warn("failed to accept connection", "error", err)
				return err
			}
		}
		go t.processConnection(conn, ctx)
	}
}

func (t *TCP) processConnection(conn net.Conn, ctx context.Context) {
	t.wg.Add(1)

	defer t.wg.Done()
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			t.logger.Warn("failed to close connection", "error", err)
		}
	}(conn)

	client := t.cm.HandleNewClient(conn)
	t.logger.Info("client: " + client.ID + " is connected")

	for {
		select {
		case <-ctx.Done():
			t.logger.Info("client connection handler stopping", "client_id", client.ID)
			return
		default:
		}
	}
}

func (t *TCP) Start(b *core.Broker, ctx context.Context) error {
	t.logger.Info("starting TCP component")

	t.broker = b
	t.cm = core.NewClientManager(b)

	t.ctx, t.cancel = context.WithCancel(ctx)

	go func() {
		t.logger.Info("TCP server running", "port", 9091)
		if err := t.HandleConnection(t.ctx); err != nil && err != context.Canceled {
			t.logger.Warn("TCP Handler failed", "error", err)
		}
	}()

	return nil
}

func (t *TCP) Stop(ctx context.Context) error {
	t.logger.Info("stopping TCP Component")

	if t.cancel != nil {
		t.cancel()
	}

	if err := t.Listener.Close(); err != nil {
		t.logger.Warn("error closing listener", "error", err)
	}

	done := make(chan struct{})
	go func() {
		t.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		t.logger.Info("All connections closed gracefully")
	case <-time.After(10 * time.Second):
		t.logger.Warn("Timeout waiting for connections to close")
		t.cm.Shutdown(ctx) // Force close
	}

	return nil
}
