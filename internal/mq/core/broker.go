// Package core provides the core components of the HopperMQ system, including the Broker component.
package core

import (
	"context"
	"log/slog"
	"sync"

	"github.com/hoppermq/hopper/internal/mq/core/protocol/container"
	"github.com/hoppermq/hopper/internal/mq/core/protocol/serializer"
	"github.com/hoppermq/hopper/pkg/domain"
)

// Broker is the core component of the HopperMQ system, responsible for managing message queues and handling client connections.
type Broker struct {
	Logger *slog.Logger

	services   []domain.Service
	Serializer *serializer.Serializer // should create a domain type here

	eb               domain.IEventBus
	cm               *ClientManager
	containerManager *container.ContainerManager

	wg     sync.WaitGroup
	cancel context.CancelFunc
}

// NewBroker creates a new Broker instance with all its core dependencies
func NewBroker(logger *slog.Logger, serializer *serializer.Serializer) *Broker {
	broker := &Broker{
		Logger:     logger,
		Serializer: serializer,
	}

	broker.cm = NewClientManager(broker)
	broker.containerManager = container.NewContainerManager()

	return broker
}

func (b *Broker) spawnHandler(ctx context.Context, eventHandler func(ctx2 context.Context)) {
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		eventHandler(ctx)
	}()
}

// Run initializes the Broker component, setting up necessary resources and preparing it to handle incoming frames and client connections.
func (b *Broker) Run(ctx context.Context, transports ...domain.Service) error {
	b.Logger.Info("Starting Broker Component")

	ctx, b.cancel = context.WithCancel(ctx)

	rcvdFrameCh := b.eb.Subscribe(string(domain.EventTypeReceiveMessage))
	newConnCh := b.eb.Subscribe(string(domain.EventTypeNewConnection))
	closedConnCh := b.eb.Subscribe(string(domain.EventTypeConnectionClosed))

	b.spawnHandler(ctx, func(ctx context.Context) {
		b.onReceivingMessage(ctx, rcvdFrameCh)
	})

	b.spawnHandler(ctx, func(ctx context.Context) {
		b.onNewClientConnection(ctx, newConnCh)
	})

	b.spawnHandler(ctx, func(ctx context.Context) {
		b.onClientDisconnect(ctx, closedConnCh)
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

	// Cancel context to stop all handlers
	if b.cancel != nil {
		b.cancel()
	}

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

func (b *Broker) RegisterEventBus(eb domain.IEventBus) {
	b.eb = eb
	b.Logger.Info("EventBus registered with", "service", b.Name())
}
