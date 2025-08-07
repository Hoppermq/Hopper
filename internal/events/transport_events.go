package events

import (
	"github.com/hoppermq/hopper/pkg/domain"
	"net"
)

type BaseEvent struct {
	EventType domain.EventType
	ClientID  string
}

func (evt *BaseEvent) GetType() domain.EventType {
	return domain.EventType(evt.EventType)
}

type NewConnectionEvent struct {
	ClientID   string
	Connection net.Conn // hope i can do this bahahha
	Transport  string

	BaseEvent
}

type MessageReceivedEvent struct {
	ClientID  string // i don't think transport is aware of this.
	Message   []byte
	Transport string

	BaseEvent
}

type SendMessageEvent struct {
	ClientID  string //should be the conn here
	Message   []byte
	Transport string

	BaseEvent
}
