package domain

import "context"

// Service represent the contract for the Service that can be used by the application or as subservice by the broker.
type Service interface {
	Name() string
	Run(ctx context.Context) error
	Stop(ctx context.Context) error
}

// EventBusAware represent the contract for service that can handle event bus channels.
type EventBusAware interface {
	RegisterEventBus(eventBus IEventBus)
}
