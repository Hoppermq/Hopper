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
func NewContainer(id, clientID domain.ID) *Container {
	return &Container{
		ID:              id,
		ClientID:        clientID,
		State:           domain.ContainerCreated,
		Channels:        make(map[domain.ID]domain.Channel),
		ChannelsByTopic: make(map[string]domain.ID),
	}
}

// CreateChannel create a new Channel and attach it to the container.
func (ctr *Container) CreateChannel(
	topic string,
	generateIdentifier func() domain.ID,
) *Channel {
	channel := NewChannel(generateIdentifier, topic)
	ctr.Channels[channel.ID] = channel
	ctr.ChannelsByTopic[topic] = channel.ID

	return channel
}

// RemoveChannel remove the channel from the container.
func (ctr *Container) RemoveChannel(topic string) {
	chanToRemove := ctr.findChannelByTopic(topic)
	if chanToRemove != nil {
		delete(ctr.Channels, chanToRemove.GetID())
		delete(ctr.ChannelsByTopic, topic)
	}
}

// SetState set the current container state.
func (ctr *Container) SetState(state domain.ContainerState) {
	ctr.State = state
}

// GetID return the containerID.
func (ctr *Container) GetID() domain.ID {
	return ctr.ID
}
