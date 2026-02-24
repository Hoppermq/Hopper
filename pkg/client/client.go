// Package client represent the hoppermq client.
package client

import (
	"context"
	"log/slog"
	"sync"

	"github.com/hoppermq/hopper/pkg/client/config"
	"github.com/hoppermq/hopper/pkg/client/transport/tcp"
	"github.com/hoppermq/hopper/pkg/domain"
)

type ClientState bool

// Client represent the sdk client.
type Client struct {
	id    domain.ID
	state ClientState

	conn      domain.Connection
	transport domain.Transport

	subscriptions     map[string]string
	subscriptionsByID map[domain.ID]string

	inboundQueue  chan string
	outboundQueue chan string

	mu     sync.RWMutex
	wg     sync.WaitGroup
	cancel context.CancelFunc

	logger *slog.Logger
}

// Option type represent the injection function.
type Option func(*Client)

func WithConfig(cfg *config.ClientConfig) Option {
	return func(c *Client) {}
}

func WithLogger(logger *slog.Logger) Option {
	return func(c *Client) {
		c.logger = logger
	}
}

func withTransport() Option {
	return func(c *Client) {
		tcpClient := tcp.NewTCPClient(
			tcp.WithLogger(c.logger),
		)
		c.transport = tcpClient
	}
}

// NewClient create a new client.
func NewClient(opts ...Option) *Client {
	c := &Client{}

	opts = append(opts, withTransport())
	for _, opts := range opts {
		opts(c)
	}

	return c
}

// Run start the client sdk workers.
func (c *Client) Run(ctx context.Context) error {
	c.logger.Info("starting hopperMQ client")
	ctx, c.cancel = context.WithCancel(ctx)

	c.setState(true)
	if err := c.transport.Run(ctx); err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}

// Stop gracefully shutdown the client sdk.
func (c *Client) Stop(ctx context.Context) error {
	c.logger.Info("stopping hopperMQ client", "client_id", c.id)
	c.setState(false)

	if c.cancel != nil {
		c.cancel()
	}

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			panic(err)
		}
	}

	c.wg.Wait()
	c.logger.Info("hopperMQ client stopped successfully")
	return nil
}

func (c *Client) setState(state ClientState) {
	c.state = state
}
