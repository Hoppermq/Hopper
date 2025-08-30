package domain

import "context"

// Transport represent the contract for our kind of Transport.
type Transport interface {
	Service
	HandleConnection(ctx context.Context) error
	RegisterEventBus(eb IEventBus)
}
