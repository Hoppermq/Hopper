package mq

import (
	"context"
	"log/slog"

	"github.com/hoppermq/hopper/internal/handler"
)

type HopperMQService struct {
	logger     *slog.Logger
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

func New(opts ...Option) *HopperMQService {
	service := &HopperMQService{}
	for _, opt := range opts {
		opt(service)
	}

	return service
}

func (h *HopperMQService) Name() string {
	return ""
}

func (h *HopperMQService) Run(ctx context.Context) error {
	return nil
}

func (h *HopperMQService) Stop(ctx context.Context) error {
	return nil
}
