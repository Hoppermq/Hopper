package events

import (
	"github.com/hoppermq/hopper/pkg/domain"
)

type BaseEvent struct {
	EventType string
}

func (evt *BaseEvent) GetType() domain.EventType {
	return domain.EventType(evt.EventType)
}

type NewConnectionEvent struct {
	Conn      domain.Connection
	Transport domain.TransportType

	BaseEvent
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
	Transport string

	BaseEvent
}

type SendMessageEvent struct {
	ClientID  domain.ID
	Conn      domain.Connection
	Message   []byte
	Transport string

	BaseEvent
}
