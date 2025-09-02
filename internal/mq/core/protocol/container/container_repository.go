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
func (mgr *Manager) FindContainersByTopic(topic string) []Container {
	mgr.Registry.mu.RLock()
	defer mgr.Registry.mu.RUnlock()

	var containers []Container

	containersID := mgr.Registry.data[topic]
	for containerID, _ := range containersID {
		container := mgr.Containers[containerID]
		if container != nil {
			containers = append(containers, *container)
		}
	}

	return containers
}
