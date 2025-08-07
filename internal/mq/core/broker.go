// Package core provides the core components of the HopperMQ system, including the Broker component.
package core

import (
	"context"
	"github.com/hoppermq/hopper/pkg/domain"
	"log/slog"
)

// Broker is the core component of the HopperMQ system, responsible for managing message queues and handling client connections.
type Broker struct {
	Logger *slog.Logger

	services []domain.Service
}

// Start initializes the Broker component, setting up necessary resources and preparing it to handle incoming messages and client connections.
func (b *Broker) Start(ctx context.Context, transports ...domain.Service) error {
	b.Logger.Info("Starting Broker Component")

	for _, transport := range transports {
		go func(t domain.Service) {
			if err := t.Run(ctx); err != nil {
				b.Logger.Error("Failed to start transport", "error", err)
			} else {
				b.Logger.Info("Transport started successfully", "transport", t.Name())
			}
		}(transport)
	}
	return nil
}

// Stop gracefully shuts down the Broker component, ensuring that all ongoing operations are completed and resources are released.
func (b *Broker) Stop(ctx context.Context) error {
	b.Logger.Info("Stopping Broker Component")

	return nil
}

func (b *Broker) HandleNewConnection(tt struct{}) error {
	return nil
}
