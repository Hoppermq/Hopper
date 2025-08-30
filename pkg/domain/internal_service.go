package domain

// IService represents the interface of our core service for the broker
type IService interface {
	Service
	RegisterEventBus(eb IEventBus)
}
