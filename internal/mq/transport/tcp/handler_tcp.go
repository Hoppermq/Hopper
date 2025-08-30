// Package transport provides all the transport handler for the Hopper message queue system.
package tcp

import (
	"bufio"
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/hoppermq/hopper/internal/events"
	"github.com/hoppermq/hopper/pkg/domain"
)

// TCP is an TCP handler.
type TCP struct {
	Listener net.Listener
	logger   *slog.Logger

	eb domain.IEventBus

	cancel context.CancelFunc
	wg     sync.WaitGroup
}

type config struct {
	lconf  *net.ListenConfig
	logger *slog.Logger
}

type Option func(*config) error

func WithLogger(logger *slog.Logger) Option {
	return func(c *config) error {
		c.logger = logger

		return nil
	}
}

func WithListener(listenerConfig *net.ListenConfig) Option {
	return func(c *config) error {
		c.lconf = listenerConfig

		return nil
	}
}

// NewTCP return the new tcp handler.
func NewTCP(ctx context.Context, opts ...Option) (*TCP, error) {
	handlerConfig := &config{}
	for _, opt := range opts {
		err := opt(handlerConfig)
		if err != nil {
			return nil, err
		}
	}

	l, err := handlerConfig.lconf.Listen(ctx, "tcp", ":9091")
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
				t.logger.Info("context canceled, stopping connection handler")
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

	if t.eb == nil {
		t.logger.Warn("EventBus not registered, skipping event publishing")
		return
	}

	evt := &events.NewConnectionEvent{
		Conn:      conn,
		Transport: domain.TransportTypeTCP,
		BaseEvent: events.BaseEvent{
			EventType: domain.EventTypeNewConnection,
		},
	}

	if err := t.eb.Publish(ctx, evt); err != nil {
		t.logger.Warn("failed to publish new connection event", "error", err)
		return
	}

	t.logger.Info(
		"channel event published",
		"publisher",
		t.Name(),
		"event",
		evt.GetType(),
		"transport",
		evt.GetType(),
	)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := t.receiveMsg(conn, ctx)
			if err != nil {
				return
			}
		}
	}
}

func (t *TCP) receiveMsg(conn domain.Connection, ctx context.Context) error {
	reader := bufio.NewReader(conn)
	if err := conn.SetReadDeadline(time.Now().Add(50 * time.Second)); err != nil {
		t.logger.Warn("failed to set read deadline", "error", err)
		return err
	}

	msg, err := reader.ReadBytes('\n')
	if err != nil {
		if err == io.EOF {
			t.logger.Info("client disconnected from tcp")
			evt := &events.ClientDisconnectedEvent{
				Transport: domain.TransportTypeTCP,
				Conn:      conn,
				BaseEvent: events.BaseEvent{
					EventType: domain.EventTypeConnectionClosed,
				},
			}

			if err := t.eb.Publish(ctx, evt); err != nil {
				t.logger.Warn("failed to publish client disconnected event", "error", err)
				return err
			}

			return err
		}
	}

	t.logger.Info("new message receive", "content", msg)
	if t.eb == nil {
		t.logger.Warn("EventBus not registered, skipping event publishing")
		return err
	}

	evt := &events.MessageReceivedEvent{
		Message:   msg,
		Transport: domain.TransportTypeTCP,
		BaseEvent: events.BaseEvent{
			EventType: domain.EventTypeReceiveMessage,
		},
	}

	if err := t.eb.Publish(ctx, evt); err != nil {
		t.logger.Warn("failed to publish message event", "error", err)
		return err
	}

	return nil
}

// Run wil start the tcp component.
func (t *TCP) Run(ctx context.Context) error {
	t.logger.Info("starting TCP component")

	ctx, t.cancel = context.WithCancel(ctx)

	msgSenderCh := t.eb.Subscribe(string(domain.EventTypeSendMessage))

	t.spawnHandler(ctx, func(ctx context.Context) {
		t.handleMessageSending(ctx, msgSenderCh)
	})

	go func() {
		t.logger.Info("TCP server running", "port", 9091)
		if err := t.HandleConnection(ctx); err != nil && !errors.Is(err, context.Canceled) {
			t.logger.Warn("TCP Handler failed", "error", err)
		}
	}()

	return nil
}

// Stop will shut down gracefully the tcp component.
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

// Name will return the component name.
func (t *TCP) Name() string {
	return "tcp-handler"
}

// RegisterEventBus will attach the event bus to the component.
func (t *TCP) RegisterEventBus(eb domain.IEventBus) {
	t.eb = eb
	t.logger.Info("EventBus registered with TCP", "service", t.Name())
}

func (t *TCP) handleMessageSending(ctx context.Context, ch <-chan domain.Event) {
	for {
		select {
		case <-ctx.Done():
			return
		case evt, ok := <-ch:
			if !ok {
				return
			}
			if c, ok := evt.(*events.SendMessageEvent); ok {
				t.logger.Info("new message event received", "transport", c.Transport, "clientID", c.ClientID)
				t.sendMessage(c.Message, c.Conn)
			}
		}
	}
}

func (t *TCP) spawnHandler(ctx context.Context, eventHandler func(ctx2 context.Context)) {
	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		eventHandler(ctx)
	}()
}
