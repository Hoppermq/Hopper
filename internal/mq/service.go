package mq

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/hoppermq/hopper/pkg/domain"

	"github.com/hoppermq/hopper/internal/common"
	"github.com/hoppermq/hopper/internal/mq/core"
	"github.com/hoppermq/hopper/internal/mq/core/protocol/serializer"
	handler "github.com/hoppermq/hopper/internal/mq/transport/tcp"
)

type HopperMQService struct {
	eb domain.IEventBus

	broker     *core.Broker
	tcpHandler *handler.TCP

	logger *slog.Logger

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

type Option func(*HopperMQService)

func WithLogger(logger *slog.Logger) Option {
	return func(s *HopperMQService) {
		s.logger = logger
	}
}

func WithTCP(tcpHandler *handler.TCP) Option {
	return func(s *HopperMQService) {
		s.tcpHandler = tcpHandler
	}
}

func WithBroker(hopperBroker *core.Broker) Option {
	return func(s *HopperMQService) {
		s.broker = hopperBroker
	}
}

func New(opts ...Option) *HopperMQService {
	service := &HopperMQService{}
	for _, opt := range opts {
		opt(service)
	}

	serializer := serializer.NewSerializer(
		common.NewPool(func() *bytes.Buffer {
			return &bytes.Buffer{}
		}),
	)

	service.broker = core.NewBroker(service.logger, serializer)

	return service
}

func (h *HopperMQService) Name() string {
	return "hopper-mq" // should be loaded from config.
}

func (h *HopperMQService) startService(name string, runner func() error) {
	h.logger.Info("Starting service", "service", name)
	if err := runner(); err != nil && !errors.Is(err, context.Canceled) {
		h.logger.Error("Service failed", "service", name, "error", err)

		h.cancel()
	}
}

func (h *HopperMQService) Run(ctx context.Context) error {
	h.ctx, h.cancel = context.WithCancel(ctx) // should be init at the main prob
	h.wg.Add(1)

	go h.startService("broker", func() error {
		defer h.wg.Done()
		return h.broker.Start(h.ctx, h.tcpHandler)
	})

	<-h.ctx.Done()
	h.logger.Info("Service stopped", "service", h.Name())

	return nil
}

func (h *HopperMQService) Stop(ctx context.Context) error {
	h.logger.Info("Stopping services")

	if h.cancel != nil {
		h.cancel()
	}

	if err := h.tcpHandler.Stop(ctx); err != nil {
		h.logger.Error("Failed to stop TCP HAndler")
	}

	if err := h.broker.Stop(ctx); err != nil {
		h.logger.Error("Failed to stop broker")
	}

	h.logger.Info("Services stopped")
	return nil
}

func (h *HopperMQService) RegisterEventBus(eb domain.IEventBus) {
	h.eb = eb

	if h.broker != nil {
		h.broker.RegisterEventBus(eb)
	}
	if h.tcpHandler != nil {
		h.tcpHandler.RegisterEventBus(eb)
	}
}
