package events

import (
	"context"
	"github.com/hoppermq/hopper/pkg/domain"
	"sync"
)

type Event struct {
	domain.Event

	Type    domain.EventType
	Payload any
}

type EventBus struct {
	mu        sync.RWMutex
	channels  map[domain.EventType][]domain.EventChannel
	maxBuffer uint16
}

func NewEventBus(maxBuffer uint16) *EventBus {
	return &EventBus{
		channels:  make(map[domain.EventType][]domain.EventChannel),
		maxBuffer: maxBuffer,
	}
}

func (eb *EventBus) getSubscribers(eventType domain.EventType) []domain.EventChannel {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	subs, ok := eb.channels[eventType]
	if !ok {
		return nil
	}

	return subs
}

func (eb *EventBus) Subscribe(eventType string) <-chan domain.Event {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	ch := make(chan domain.Event, eb.maxBuffer)
	eb.channels[domain.EventType(eventType)] = append(eb.channels[domain.EventType(eventType)], ch)

	return ch
}

func (eb *EventBus) Publish(ctx context.Context, event domain.Event) error {
	eb.mu.RLock()
	subs := eb.getSubscribers(event.Type())
	eb.mu.RUnlock()

	if len(subs) == 0 {
		return nil
	}

	var dropped int

	for _, sub := range subs {
		select {
		case sub <- event:
		default:
			dropped++
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	if dropped > 0 {
		return domain.ErrInvalidHeader
	}

	return nil
}
