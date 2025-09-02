package container

import "github.com/hoppermq/hopper/pkg/domain"

func (ctr *Container) findChannelByTopic(topic string) domain.Channel {
	if channelID, ok := ctr.ChannelsByTopic[topic]; ok {
		return ctr.Channels[channelID]
	}
	return nil
}

func (ctr *Container) findChannelByID(id domain.ID) domain.Channel {
	return ctr.Channels[id]
}
