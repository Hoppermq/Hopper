package events

import (
	"context"
	"sync"

	"github.com/hoppermq/hopper/internal/config"
	"github.com/hoppermq/hopper/pkg/domain"
)

type EventBus struct {
	mu sync.RWMutex
	wg sync.WaitGroup
	configuration *config.Configuration

	channels  map[domain.EventType][]domain.EventChannel
	maxBuffer uint16
}

func NewEventBus(maxBuffer uint16) *EventBus {
	return &EventBus{
		channels:  make(map[domain.EventType][]domain.EventChannel),
		maxBuffer: maxBuffer,
	}
}

// Option type represent the options of the event bus.
type Option func(*EventBus)

// WithConfig inject the configuration to the event bus.
func WithConfig(cfg *config.Configuration) Option {
	return func(e *EventBus) {
		println("config:", cfg)
		e.configuration = cfg
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
	subs := eb.getSubscribers(event.GetType())

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
