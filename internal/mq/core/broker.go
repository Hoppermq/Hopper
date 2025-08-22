// Package core provides the core components of the HopperMQ system, including the Broker component.
package core

import (
	"context"
	"log/slog"
	"sync"

	"github.com/hoppermq/hopper/internal/events"
	"github.com/hoppermq/hopper/internal/mq/core/protocol/container"
	"github.com/hoppermq/hopper/internal/mq/core/protocol/frames"
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

	// Broker manages its own core dependencies
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

// Start initializes the Broker component, setting up necessary resources and preparing it to handle incoming messages and client connections.
func (b *Broker) Start(ctx context.Context, transports ...domain.Service) error {
	b.Logger.Info("Starting Broker Component")

	ctx, b.cancel = context.WithCancel(ctx)

	connCh := b.eb.Subscribe(string(domain.EventTypeNewConnection))

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

func (b *Broker) RegisterEventBus(eb domain.IEventBus) {
	b.eb = eb
	b.Logger.Info("EventBus registered with", "service", b.Name())
}

func (b *Broker) onNewClientConnection(ctx context.Context, evt *events.NewConnectionEvent) {
	client := b.cm.HandleNewClient(evt.Conn)
	b.Logger.Info("New client connection handled", "clientID", client.ID)

	openFramePayloadData := frames.CreateOpenFramePayloadData(client.ID, GenerateIdentifier())
	data, err := b.Serializer.SerializeOpenFramePayloadData(openFramePayloadData)
	if err != nil {
		b.Logger.Warn("failed to serialize payload data")
		return
	}

	openFrame, err := frames.CreateOpenFrame(domain.DOFF2, data)
	if err != nil {
		b.Logger.Warn("failed to create open frame", "error", err)
		return
	}

	frame, err := b.Serializer.SerializeFrame(openFrame)

	if err != nil {
		b.Logger.Error("failed to serialize open frame", "error", err)
		return
	}

	sendMsgEvt := &events.SendMessageEvent{
		ClientID:  client.ID,
		Conn:      client.Conn,
		Message:   frame,
		Transport: string(domain.TransportTypeTCP),
		BaseEvent: events.BaseEvent{
			EventType: string(domain.EventTypeSendMessage),
		},
	}

	if err := b.eb.Publish(ctx, sendMsgEvt); err != nil {
		b.Logger.Error("Failed to publish SendMessageEvent", "error", err)
		return
	}

	b.Logger.Info("SendMessageEvent published", "clientID", client.ID, "transport", sendMsgEvt.Transport, "message", string(sendMsgEvt.Message))
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
				b.onNewClientConnection(ctx, c)
			}
		}
	}
}
