// Package transport provides all the transport handler for the Hopper message queue system.
package transport

import (
	"context"
	"errors"
	"github.com/hoppermq/hopper/internal/events"
	"github.com/hoppermq/hopper/pkg/domain"
	"log/slog"
	"net"
	"sync"
	"time"
)

// TCP is an TCP handler
type TCP struct {
	Listener net.Listener
	logger   *slog.Logger

	eb domain.IEventBus

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
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

func (t *TCP) processConnection(conn domain.Connection, ctx context.Context) {
	t.wg.Add(1)
	defer t.wg.Done()

	defer func(conn domain.Connection) {
		err := conn.Close()
		if err != nil {
			t.logger.Warn("failed to close connection", "error", err)
		}
	}(conn)

	// prob a goroutine to send events to the event bus
	if t.eb == nil {
		t.logger.Warn("EventBus not registered, skipping event publishing")
		return
	}

	evt := &events.NewConnectionEvent{
		Conn:      conn,
		Transport: string(domain.TransportTypeTCP),
		BaseEvent: events.BaseEvent{
			EventType: string(domain.EventTypeNewConnection),
		},
	}

	if err := t.eb.Publish(ctx, evt); err != nil {
		t.logger.Warn("failed to publish new connection event", "error", err)
		return
	}
	t.logger.Info("channel event published", "publisher", t.Name(), "event", evt.EventType, "transport", evt.Transport)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

func (t *TCP) Run(ctx context.Context) error {
	t.logger.Info("starting TCP component")

	t.ctx, t.cancel = context.WithCancel(ctx)

	go func() {
		t.logger.Info("TCP server running", "port", 9091)
		if err := t.HandleConnection(t.ctx); err != nil && !errors.Is(err, context.Canceled) {
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
	}

	return nil
}

func (t *TCP) Name() string {
	return "tcp-handler"
}

func (t *TCP) RegisterEventBus(eb *events.EventBus) {
	t.eb = eb
	t.logger.Info("EventBus registered with TCP", "service", t.Name())
}
