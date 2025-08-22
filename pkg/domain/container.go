package domain

type Container interface {
	CreateChannel(topic string, idGenerator func() string) Channel
}

type Channel interface {
}
