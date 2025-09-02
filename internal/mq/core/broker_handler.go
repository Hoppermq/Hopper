package core

import (
	"context"

	"github.com/hoppermq/hopper/internal/common"
	"github.com/hoppermq/hopper/internal/events"
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
			Size:    10, //idk
			Type:    domain.FrameTypeOpen,
			DOFF:    domain.DOFF4,
			Channel: 0, // ?
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
	}
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
