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
			// Handle both event types for backward compatibility
			if c, ok := evt.(*events.ClientDisconnectEvent); ok {
				b.handleConnectionClosed(ctx, c)
			} else if c, ok := evt.(*events.ClientDisconnectedEvent); ok {
				b.handleConnectionClosedByConn(ctx, c)
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
				frameType := frame.GetType()
				switch {
				case b.fm.IsMessageFrame(frameType):
					b.Logger.Info("message frame received", "frame_type", frameType)
				case b.fm.IsControlFrame(frameType):
					b.Logger.Info("control frame received", "frame_type", frameType)
					b.RouteControlFrames(frame)
				case b.fm.IsErrorFrame(frameType):
					b.Logger.Info("error frame received", "frame_type", frameType)
				}
			}
		}
	}
}
