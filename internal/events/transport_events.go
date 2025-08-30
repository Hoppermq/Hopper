// Package events represent the channel event bus.
package events

import (
	"github.com/hoppermq/hopper/pkg/domain"
)

// BaseEvent will handle connection and transport or at least transport.
type BaseEvent struct {
	EventType domain.EventType
}

// GetType return the event type.
func (evt *BaseEvent) GetType() domain.EventType {
	return evt.EventType
}

// NewConnectionEvent represent the event for a new connection.
type NewConnectionEvent struct {
	Conn      domain.Connection
	Transport domain.TransportType

	BaseEvent
}

// GetType return the eventType.
func (evt *NewConnectionEvent) GetType() domain.EventType {
	return evt.EventType
}

// GetTransport return the transport used.
func (evt *NewConnectionEvent) GetTransport() domain.TransportType {
	return evt.Transport
}

// ClientDisconnectEvent represent the disconnection event of a client.
type ClientDisconnectEvent struct {
	ClientID  domain.ID
	Transport domain.TransportType

	Conn domain.Connection

	BaseEvent
}

// GetTransport return the transport used.
func (evt *ClientDisconnectEvent) GetTransport() domain.TransportType {
	return evt.Transport
}

// MessageReceivedEvent represent a new message received event.
type MessageReceivedEvent struct {
	Message   []byte
	Transport domain.TransportType

	BaseEvent
}

// GetType return the eventType.
func (evt *MessageReceivedEvent) GetType() domain.EventType {
	return evt.EventType
}

// GetTransport return  the transport used.
func (evt *MessageReceivedEvent) GetTransport() domain.TransportType {
	return evt.Transport
}

// SendMessageEvent represent a new message sent evet.
type SendMessageEvent struct {
	ClientID  domain.ID
	Conn      domain.Connection
	Message   []byte
	Transport domain.TransportType

	BaseEvent
}

// GetType return the eventType.
func (evt *SendMessageEvent) GetType() domain.EventType {
	return evt.EventType
}

// GetTransport return the transport used.
func (evt *SendMessageEvent) GetTransport() domain.TransportType {
	return evt.Transport
}

// ClientDisconnectedEvent represent th event of client disconnection.
type ClientDisconnectedEvent struct {
	Conn      domain.Connection
	Transport domain.TransportType

	BaseEvent
}

// GetType return the eventType.
func (evt *ClientDisconnectedEvent) GetType() domain.EventType {
	return evt.EventType
}

// GetTransport return the transport used.
func (evt *ClientDisconnectedEvent) GetTransport() domain.TransportType {
	return evt.Transport
}
