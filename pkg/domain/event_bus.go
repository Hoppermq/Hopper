package domain

type EventChannel chan Event

type EventType string

type Event interface {
	Type() EventType // will be typed later
}

type EventBus interface {
	Publish(event Event) error
	Subscribe(eventType string, handler func(event Event)) error
}
