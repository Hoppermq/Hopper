package core

import (
	"context"

	"github.com/hoppermq/hopper/internal/events"
	"github.com/hoppermq/hopper/pkg/domain"
)

func (b *Broker) onNewClientConnection(ctx context.Context, ch <-chan domain.Event) {
	for {
		select {
		case <-ctx.Done():
			return
		case evt, ok := <-ch:
			if !ok {
				return
			}
			if c, ok := evt.(*events.NewConnectionEvent); ok {
				b.handleNewClientConnection(ctx, c)
			}
		}
	}
}

func (b *Broker) onClientDisconnect(ctx context.Context, ch <-chan domain.Event) {
	for {
		select {
		case <-ctx.Done():
			return
		case evt, ok := <-ch:
			if !ok {
				return
			}
			if c, ok := evt.(*events.ClientDisconnectEvent); ok {
				b.handleConnectionClosed(ctx, c)
			}
		}
	}
}

func (b *Broker) onReceivingMessage(ctx context.Context, ch <-chan domain.Event) {
	for {
		select {
		case <-ctx.Done():
			return
		case evt, ok := <-ch:
			if !ok {
				return
			}
			if c, ok := evt.(*events.MessageReceivedEvent); ok {
				frame, err := b.Serializer.DeserializeFrame(c.Message)
				if err != nil {
					b.Logger.Warn("failed to deserialize the frame", "error", err)
					continue
				}

				b.Logger.Info("new frame received", "frame_type", frame.GetType())
			}
		}
	}
}
