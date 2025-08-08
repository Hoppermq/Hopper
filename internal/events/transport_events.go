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
	Transport string

	BaseEvent
}

func (evt *NewConnectionEvent) GetTransport() domain.TransportType {
	return domain.TransportType(evt.Transport)
}

type MessageReceivedEvent struct {
	Message   []byte
	Transport string

	BaseEvent
}

type SendMessageEvent struct {
	Message   []byte
	Transport string

	BaseEvent
}
