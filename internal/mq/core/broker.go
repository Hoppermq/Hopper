// Package core provides the core components of the HopperMQ system, including the Broker component.
package core

import (
	"context"
	"github.com/hoppermq/hopper/internal/events"
	"github.com/hoppermq/hopper/pkg/domain"
	"log/slog"
	"sync"
)

// Broker is the core component of the HopperMQ system, responsible for managing message queues and handling client connections.
type Broker struct {
	Logger *slog.Logger

	services []domain.Service
	eb       *events.EventBus
	cm       *ClientManager

	wg     sync.WaitGroup
	cancel context.CancelFunc
}

func (b *Broker) spawnHandler(ctx context.Context, eventHandler func(ctx2 context.Context)) {
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		eventHandler(ctx)
	}()
}

// Start initializes the Broker component, setting up necessary resources and preparing it to handle incoming messages and client connections.
func (b *Broker) Start(ctx context.Context, transports ...domain.Service) error {
	b.Logger.Info("Starting Broker Component")

	ctx, b.cancel = context.WithCancel(ctx)

	connCh := b.eb.Subscribe("new_connection")

	b.spawnHandler(ctx, func(ctx context.Context) {
		b.handleNewConnections(ctx, connCh)
	})

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
	for _, service := range b.services {
		if err := service.Stop(ctx); err != nil {
			b.Logger.Error("Failed to stop service", "service", service.Name(), "error", err)
		} else {
			b.Logger.Info("Service stopped successfully", "service", service.Name())
		}
	}

	return nil
}

func (b *Broker) Name() string {
	return "hopper-broker"
}

func (b *Broker) RegisterEventBus(eb *events.EventBus) {
	b.eb = eb
	b.Logger.Info("EventBus registered with", "service", b.Name())
}

func (b *Broker) onNewClientConnection(evt *events.NewConnectionEvent) {
	client := b.cm.HandleNewClient(evt.Conn)
	b.Logger.Info("New client connection handled", "clientID", client.ID)
	// should publish frame here
}

func (b *Broker) handleNewConnections(ctx context.Context, ch <-chan domain.Event) {
	for {
		select {
		case <-ctx.Done():
			return
		case evt, ok := <-ch:
			if !ok {
				return
			}
			if c, ok := evt.(*events.NewConnectionEvent); ok {
				b.Logger.Info("New connection event received", "transport", c.Transport)
				b.onNewClientConnection(c)
			}
		}
	}
}
