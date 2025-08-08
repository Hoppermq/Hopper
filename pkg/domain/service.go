package domain

import "context"

type Service interface {
	Name() string
	Run(ctx context.Context) error
	Stop(ctx context.Context) error
}

type EventBusAware interface {
	RegisterEventBus(eventBus IEventBus) // You'll need EventBus interface in domain
}
