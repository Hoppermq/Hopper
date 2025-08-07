package domain

import "context"

type EventChannel chan Event

type EventType string

type Event interface {
	GetType() EventType // will be typed later
}

type IEventBus interface {
	Publish(ctx context.Context, event Event) error
	Subscribe(eventType string) <-chan Event
}
