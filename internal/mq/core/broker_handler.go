package core

import (
	"context"
	"fmt"

	"github.com/hoppermq/hopper/internal/common"
	"github.com/hoppermq/hopper/internal/events"
	"github.com/hoppermq/hopper/internal/mq/core/protocol/container"
	"github.com/hoppermq/hopper/internal/mq/core/protocol/frames"
	"github.com/hoppermq/hopper/pkg/domain"
)

func (b *Broker) handleNewClientConnection(ctx context.Context, evt *events.NewConnectionEvent) {
	client := b.clientManager.HandleNewClient(evt.Conn)
	ctr := b.containerManager.CreateNewContainer(
		common.GenerateIdentifier,
		client.ID,
	)
	b.Logger.Info(
		"new container created",
		"container_id",
		ctr.GetID(),
		"current_state",
		ctr.State,
	)
	frameHeaderPayload := &frames.PayloadHeader{}
	framePayload := frames.CreateOpenFramePayload(
		frameHeaderPayload,
		client.ID,
		ctr.GetID(),
	)

	frame, err := frames.CreateFrame(
		&frames.Header{
			Size:    10,
			Type:    domain.FrameTypeOpen,
			DOFF:    domain.DOFF4,
			Channel: 0,
		},
		nil,
		framePayload,
	)
	if err != nil {
		b.Logger.Warn("failed to create frame from payload", "error", err)
		return
	}

	data, err := b.Serializer.SerializeFrame(frame)
	if err != nil {
		b.Logger.Warn("failed to serialize frame", "error", err)
	}
	sendMsgEvt := &events.SendMessageEvent{
		ClientID:  client.ID,
		Conn:      client.Conn,
		Message:   data,
		Transport: domain.TransportTypeTCP,
		BaseEvent: events.BaseEvent{
			EventType: domain.EventTypeSendMessage,
		},
	}

	if err := b.eb.Publish(ctx, sendMsgEvt); err != nil {
		b.Logger.Warn("failed to publish new message send event", "error", err)
		return
	}

	b.containerManager.UpdateContainerState(ctr.ID, domain.ContainerOpenSent)

}

func (b *Broker) handleConnectionClosed(ctx context.Context, evt *events.ClientDisconnectEvent) {
	b.Logger.Info("client disconnected event", "client", evt.ClientID)

	b.clientManager.RemoveClient(evt.ClientID)
}

func (b *Broker) handleConnectionClosedByConn(ctx context.Context, evt *events.ClientDisconnectedEvent) {
	client := b.clientManager.GetClientByConnection(evt.Conn)
	if client == nil {
		b.Logger.Warn("client not found for disconnected connection")
		return
	}

	b.Logger.Info("client disconnected event", "client", client.ID)

	b.clientManager.RemoveClient(client.ID)
}

func (b *Broker) RouteControlFrames(ctx context.Context, frame domain.Frame) {
	frameType := frame.GetType()

	container := b.getContainerForFrame(frame)
	if container == nil {
		b.Logger.Warn("container not found for frame", "frame_type", frameType)
		return
	}

	sendCallback := b.createFrameSendCallback()

	if err := container.HandleFrame(ctx, frame, sendCallback); err != nil {
		b.Logger.Error("failed to handle frame in container",
			"frame_type", frameType,
			"container_id", container.GetID(),
			"error", err)
		return
	}

	b.Logger.Info("frame handled successfully",
		"frame_type", frameType,
		"container_id", container.GetID(),
		"container_state", container.GetState())
}

func (b *Broker) RouteMessageFrames(frame domain.Frame) {
	framePayload := frame.GetPayload().(*frames.MessageFramePayload)
	containers := b.containerManager.FindContainersByTopic(
		framePayload.GetTopic(),
	)
	for _, container := range containers {
		container.CreateChannel(
			framePayload.GetTopic(),
			common.GenerateIdentifier,
		)

	}
}

func (b *Broker) RouteErrorFrames(frame domain.Frame) {}

// getContainerForFrame extracts the container from a frame based on its type and payload
func (b *Broker) getContainerForFrame(frame domain.Frame) *container.Container {
	var sourceID domain.ID

	switch frame.GetType() {
	case domain.FrameTypeOpenRcvd:
		if payload, ok := frame.GetPayload().(domain.OpenRcvdFramePayload); ok {
			sourceID = payload.GetSourceID()
		}
	case domain.FrameTypeConnect:
		if payload, ok := frame.GetPayload().(domain.ConnectFramePayload); ok {
			sourceID = payload.GetSourceID()
		}
	case domain.FrameTypeSubscribe:
		if _, ok := frame.GetPayload().(domain.SubscribeFramePayload); ok {
			b.Logger.Warn("Subscribe frame handling not yet implemented")
			return nil
		}
	default:
		b.Logger.Warn("unsupported frame type for container lookup", "frame_type", frame.GetType())
		return nil
	}

	if sourceID == "" {
		b.Logger.Warn("could not extract source ID from frame", "frame_type", frame.GetType())
		return nil
	}

	client := b.clientManager.GetClient(sourceID)
	if client == nil {
		b.Logger.Warn("client not found", "source_id", sourceID)
		return nil
	}

	containerID := client.GetContainer()
	return b.containerManager.FindContainer(containerID)
}

func (b *Broker) createFrameSendCallback() func(context.Context, domain.Frame, domain.ID) error {
	return func(ctx context.Context, frame domain.Frame, clientID domain.ID) error {
		// Get client connection
		client := b.clientManager.GetClient(clientID)
		if client == nil {
			return fmt.Errorf("client not found: %s", clientID)
		}

		// Serialize frame
		data, err := b.Serializer.SerializeFrame(frame)
		if err != nil {
			return fmt.Errorf("failed to serialize frame: %w", err)
		}

		// Create send message event
		sendMsgEvt := &events.SendMessageEvent{
			ClientID:  clientID,
			Conn:      client.Conn,
			Message:   data,
			Transport: domain.TransportTypeTCP,
			BaseEvent: events.BaseEvent{
				EventType: domain.EventTypeSendMessage,
			},
		}

		if err := b.eb.Publish(ctx, sendMsgEvt); err != nil {
			return fmt.Errorf("failed to publish send message event: %w", err)
		}

		b.Logger.Info("frame sent via callback",
			"frame_type", frame.GetType(),
			"client_id", clientID)

		return nil
	}
}
