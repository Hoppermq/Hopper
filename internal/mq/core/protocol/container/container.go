package container

import (
	"github.com/hoppermq/hopper/pkg/domain"
)

type Container struct {
	ID       string
	ClientID string
	State    string

	Channels        map[string]domain.Channel // Channels by uuid
	ChannelsByTopic map[string]string         // storing uuid channel by topic
}

type Channel struct {
	ID string

	Topic      string
	RoutingKey string
}

func NewChannel(generator func() string, topic string) domain.Channel {
	return &Channel{
		ID:         generator(),
		Topic:      topic,
		RoutingKey: "chanID-topic-version", // i guess version should be usefull here no?
	}
}

func (ctnr *Container) CreateChannel(topic string, idGenerator func() string) domain.Channel {
	channel := NewChannel(idGenerator, topic)
	ctnr.Channels[channel.(*Channel).ID] = channel
	ctnr.ChannelsByTopic[topic] = channel.(*Channel).ID

	return channel
}
