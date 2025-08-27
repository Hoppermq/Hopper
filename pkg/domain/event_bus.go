package domain

import "context"

// EventChannel represent the type handled by the bus channels.
type EventChannel chan Event

// EventType represent the type of event.
type EventType string

// TransportType represent the type of transport that will handle a command.
type TransportType string

const (
	// EventTypeNewConnection is the type for new connections event.
	EventTypeNewConnection    EventType = "new_connection"
	EventTypeConnectionClosed EventType = "close_connection"

	// EventTypeSendMessage is the type for sending a msg to a transporter.
	EventTypeSendMessage    EventType = "send_message"
	EventTypeReceiveMessage EventType = "receive_message"
)

const (
	// TransportTypeTCP is the type for TCP tranpoter.
	TransportTypeTCP TransportType = "tcp"
)

// Event represent the event happening.
type Event interface {
	GetType() EventType
}

// IEventBus represent the type of bus event.
type IEventBus interface {
	Publish(ctx context.Context, event Event) error
	Subscribe(eventType string) <-chan Event
}
