package core

import (
	"context"

	"github.com/hoppermq/hopper/internal/events"
	"github.com/hoppermq/hopper/internal/mq/core/protocol/container"
	"github.com/hoppermq/hopper/internal/mq/core/protocol/frames"
	"github.com/hoppermq/hopper/pkg/domain"
)

func (b *Broker) handleNewClientConnection(ctx context.Context, evt *events.NewConnectionEvent) {
	client := b.cm.HandleNewClient(evt.Conn)
	ctnr := b.containerManager.CreateNewContainer(
		GenerateIdentifier,
		client.ID,
	)
	b.Logger.Info(
		"new container created",
		"container_id",
		ctnr.GetID(),
		"current_state",
		ctnr.(*container.Container).State,
	)
	frameHeaderPayload := &frames.PayloadHeader{}
	framePayload := frames.CreateOpenFramePayload(
		frameHeaderPayload,
		client.ID,
		ctnr.GetID(),
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
		Transport: string(domain.TransportTypeTCP),
		BaseEvent: events.BaseEvent{
			EventType: string(domain.EventTypeSendMessage),
		},
	}

	if err := b.eb.Publish(ctx, sendMsgEvt); err != nil {
		b.Logger.Warn("failed to publish new message send event", "error", err)
	}
}

func (b *Broker) handleConnectionClosed(ctx context.Context, evt *events.ClientDisconnectEvent) {
	b.Logger.Info("client disconnected event", "client", evt.ClientID)
}
