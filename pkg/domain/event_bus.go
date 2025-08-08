package domain

import "context"

type EventChannel chan Event

type EventType string
type TransportType string

const (
	EventTypeNewConnection EventType = "new_connection"
	EventTypeSendMessage   EventType = "send_message"
)

const (
	TransportTypeTCP TransportType = "tcp"
)

type Event interface {
	GetType() EventType // will be typed later
}

type IEventBus interface {
	Publish(ctx context.Context, event Event) error
	Subscribe(eventType string) <-chan Event
}
