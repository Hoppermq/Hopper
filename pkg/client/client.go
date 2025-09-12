// Package client represent the hoppermq client.
package client

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/hoppermq/hopper/pkg/client/config"
	"github.com/hoppermq/hopper/pkg/client/transport/tcp"
	"github.com/hoppermq/hopper/pkg/domain"
	"github.com/zixyos/glog"
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

	logger slog.Logger
}

// Run start the client sdk workers.
func (c *Client) Run(ctx context.Context) error {
	c.logger.Info("starting hopperMQ client")
	ctx, c.cancel = context.WithCancel(ctx)

	c.setState(true)
	c.transport.Run(ctx)

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

// NewClient create a new client.
func NewClient() *Client {
	ctx := context.Background()
	conf, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	logger, err := glog.New(
		glog.WithLevel(slog.LevelDebug),
		glog.WithStyle(glog.WithErrorStyle()),
		glog.WithJsonFormat(),
	)
	if err != nil {
		panic(err)
	}

	logger.Info("config loaded", "configuration", conf)

	tcpClient := tcp.NewTCPClient(
		tcp.WithLogger(logger),
		tcp.WithAddress("0.0.0.0", 5672),
		tcp.WithHealthConfig(30*time.Second, 5, 60*time.Second),
	)

	if err := tcpClient.Start(ctx); err != nil {
		logger.Error("client failed", "error", err)
	}

	if status := tcpClient.IsHealthy(); !status {
		logger.Warn("client is unhealthy", "status", "unhealthy")
	}

	return &Client{}
}
