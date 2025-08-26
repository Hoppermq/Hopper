package events

import (
	"github.com/hoppermq/hopper/pkg/domain"
)

// Base Event will handle connection and transport or at least transport
type BaseEvent struct {
	EventType domain.EventType
}

func (evt *BaseEvent) GetType() domain.EventType {
	return evt.EventType
}

type NewConnectionEvent struct {
	Conn      domain.Connection
	Transport domain.TransportType

	BaseEvent
}

func (evt *NewConnectionEvent) GetType() domain.EventType {
	return evt.EventType
}

func (evt *NewConnectionEvent) GetTransport() domain.TransportType {
	return evt.Transport
}

type ClientDisconnectEvent struct {
	ClientID  domain.ID
	Transport domain.TransportType

	Conn domain.Connection

	BaseEvent
}

func (evt *ClientDisconnectEvent) GetTransport() domain.TransportType {
	return evt.Transport
}

type MessageReceivedEvent struct {
	Message   []byte
	Transport domain.TransportType

	BaseEvent
}

func (evt *MessageReceivedEvent) GetType() domain.EventType {
	return evt.EventType
}

func (evt *MessageReceivedEvent) GetTransport() domain.TransportType {
	return evt.Transport
}

type SendMessageEvent struct {
	ClientID  domain.ID
	Conn      domain.Connection
	Message   []byte
	Transport domain.TransportType

	BaseEvent
}

func (evt *SendMessageEvent) GetType() domain.EventType {
	return evt.EventType
}

func (evt *SendMessageEvent) GetTransport() domain.TransportType {
	return evt.Transport
}

type ClientDisconnectedEvent struct {
	Conn      domain.Connection
	Transport domain.TransportType

	BaseEvent
}

func (evt *ClientDisconnectedEvent) GetType() domain.EventType {
	return evt.EventType
}

func (evt *ClientDisconnectedEvent) GetTransport() domain.TransportType {
	return evt.Transport
}
