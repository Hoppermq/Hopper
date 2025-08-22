package container

import (
	"github.com/hoppermq/hopper/pkg/domain"
)

// Container represent the struct of a Container that will handle channels.
type Container struct {
	ID       domain.ID
	ClientID domain.ID
	State    domain.ContainerState

	Channels        map[domain.ID]domain.Channel // Channels by uuid
	ChannelsByTopic map[string]domain.ID         // storing uuid channel by topic
}

// Channel represent the data struct of a Channel that will manage routing.
type Channel struct {
	ID domain.ID

	Topic      string
	RoutingKey string
}

// NewChannel create a new Channel.
func NewChannel(generator func() domain.ID, topic string) domain.Channel {
	return &Channel{
		ID:         generator(),
		Topic:      topic,
		RoutingKey: "chanID-topic-version", // i guess version should be usefull here no?
	}
}

func NewContainer(ID, clientID domain.ID) domain.Container {
	return &Container{
		ID:              ID,
		ClientID:        clientID,
		State:           domain.CONTAINER_CREATED,
		Channels:        make(map[domain.ID]domain.Channel),
		ChannelsByTopic: make(map[string]domain.ID),
	}
}

// CreateChannel create a new Channel and attach it to the container.
func (ctnr *Container) CreateChannel(
	topic string,
	generateIdentifier func() domain.ID,
) domain.Channel {
	channel := NewChannel(generateIdentifier, topic)
	ctnr.Channels[domain.ID(channel.(*Channel).ID)] = channel
	ctnr.ChannelsByTopic[topic] = channel.(*Channel).ID

	return channel
}

// RemoveChannel remove the channel from the container.
func (ctnr *Container) RemoveChannel(topic string) {
	chanToRemove := ctnr.findChannelByTopic(topic)
	// TODO: should avoid type assertion here done it before jut look
	delete(ctnr.Channels, chanToRemove.(*Channel).ID)
}

// TODO: Move to repository.
func (ctnr *Container) findChannelByTopic(topic string) domain.Channel {
	return ctnr.ChannelsByTopic[topic] // should add some validation here
}

// TODO: Move to repository.
func (ctnr *Container) findChannelByID(ID domain.ID) domain.Channel {
	return ctnr.Channels[ID]
}

func (ctnr *Container) SetState(state domain.ContainerState) {
	ctnr.State = state
}

func (ctnr *Container) GetID() domain.ID {
	return ctnr.ID
}
