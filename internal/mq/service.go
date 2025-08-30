// Package mq is the message broker package.
package mq

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/hoppermq/hopper/pkg/domain"
)

// HopperMQService represent the service orchestrator of the application.
type HopperMQService struct {
	eb domain.IEventBus

	broker     domain.IService
	tcpHandler domain.Transport

	logger *slog.Logger

	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// Option represent the type options.
type Option func(*HopperMQService)

// WithLogger inject the logger.
func WithLogger(logger *slog.Logger) Option {
	return func(s *HopperMQService) {
		s.logger = logger
	}
}

// WithBroker inject broker service.
func WithBroker(broker domain.IService) Option {
	return func(s *HopperMQService) {
		s.broker = broker
	}
}

// WithTransport inject transport layer.
func WithTransport(transport domain.Transport) Option {
	return func(s *HopperMQService) {
		s.tcpHandler = transport
	}
}

// WithEventBus inject event bus.
func WithEventBus(eventBus domain.IEventBus) Option {
	return func(s *HopperMQService) {
		s.eb = eventBus
	}
}

// New create a new Service orchestrator.
func New(opts ...Option) *HopperMQService {
	service := &HopperMQService{}
	for _, opt := range opts {
		opt(service)
	}

	return service
}

// Name return the service name.
func (h *HopperMQService) Name() string {
	return "hopper-mq" // should be loaded from config.
}

func (h *HopperMQService) startService(name string, runner func() error) {
	h.logger.Info("Starting service", "service", name)
	if err := runner(); err != nil && !errors.Is(err, context.Canceled) {
		h.logger.Error("service failed to startup", "service", name, "error", err)

		h.cancel()
	}
}

// Run will start all components.
func (h *HopperMQService) Run(ctx context.Context) error {
	ctx, h.cancel = context.WithCancel(ctx) // should be init at the main prob
	h.wg.Add(1)

	if h.eb == nil {
		h.logger.Warn("no service bus available")
		return domain.ErrNoServiceAvailable
	}

	go h.startService("broker", func() error {
		defer h.wg.Done()
		if err := h.broker.Run(ctx); err != nil {
			h.logger.Warn("failed to run broker", "error", err)
			return err
		}
		return nil
	})

	<-ctx.Done()
	h.logger.Info("Service stopped", "service", h.Name())

	return nil
}

// Stop will shut down gracefully all components.
func (h *HopperMQService) Stop(ctx context.Context) error {
	h.logger.Info("Stopping services")

	if h.cancel != nil {
		h.cancel()
	}

	if h.tcpHandler != nil {
		if err := h.tcpHandler.Stop(ctx); err != nil {
			h.logger.Error("Failed to stop TCP HAndler")
		}
	}

	if h.broker != nil {
		if err := h.broker.Stop(ctx); err != nil {
			h.logger.Error("Failed to stop broker")
		}
	}

	h.logger.Info("Services stopped")
	return nil
}

// RegisterEventBus register the event bus to the services.
func (h *HopperMQService) RegisterEventBus(bus domain.IEventBus) {
	h.eb = bus

	if h.broker != nil {
		h.broker.RegisterEventBus(bus)
	}
	if h.tcpHandler != nil {
		h.tcpHandler.RegisterEventBus(bus)
	}
}
