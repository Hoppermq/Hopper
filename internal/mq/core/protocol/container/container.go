// Package container represent the business logic of a container.
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

// GetID returns the channel ID - implements domain.Channel interface.
func (c *Channel) GetID() domain.ID {
	return c.ID
}

// NewChannel create a new Channel.
func NewChannel(generator func() domain.ID, topic string) *Channel {
	return &Channel{
		ID:         generator(),
		Topic:      topic,
		RoutingKey: "chanID-topic-version", // i guess version should be useful here no?
	}
}

// NewContainer return a new container.
func NewContainer(id, clientID domain.ID) domain.Container {
	return &Container{
		ID:              id,
		ClientID:        clientID,
		State:           domain.ContainerCreated,
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
	ctnr.Channels[domain.ID(channel.ID)] = channel
	ctnr.ChannelsByTopic[topic] = channel.ID

	return channel
}

// RemoveChannel remove the channel from the container.
func (ctnr *Container) RemoveChannel(topic string) {
	chanToRemove := ctnr.findChannelByTopic(topic)
	if chanToRemove != nil {
		delete(ctnr.Channels, chanToRemove.GetID())
		delete(ctnr.ChannelsByTopic, topic)
	}
}

// TODO: Move to repository.
func (ctnr *Container) findChannelByTopic(topic string) domain.Channel {
	if channelID, ok := ctnr.ChannelsByTopic[topic]; ok {
		return ctnr.Channels[channelID]
	}
	return nil
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
