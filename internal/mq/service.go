package mq

import (
	"context"
	"log/slog"

	"github.com/hoppermq/hopper/internal/handler"
	"github.com/hoppermq/hopper/internal/mq/core"
)

type HopperMQService struct {
	logger *slog.Logger
	// core logic
	broker *core.Broker
	// handler  called by the logic
	tcpHandler *handler.TCP
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

	service.broker = &core.Broker{
		Logger: service.logger,
	}

	return service
}

func (h *HopperMQService) Name() string {
	return ""
}

func (h *HopperMQService) Run(ctx context.Context) error {
	h.broker.Start()

	// goroutine for sub-services with inj of broker here
	return nil
}

func (h *HopperMQService) Stop(ctx context.Context) error {
	return nil
}
