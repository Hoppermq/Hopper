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

// FindContainersByTopic return all container attached to a topic as subscriber.
func (cm *Manager) FindContainersByTopic(topic string) []Container {
	cm.Registry.mu.RLock()
	defer cm.Registry.mu.RUnlock()

	var containers []Container

	containersID := cm.Registry.data[topic]
	for containerID, _ := range containersID {
		container := cm.Containers[containerID]
		if container != nil {
			containers = append(containers, *container)
		}
	}

	return containers
}
