package container

import "github.com/hoppermq/hopper/pkg/domain"

// FindContainerByClientID return the given container.
func (mgr *Manager) FindContainerByClientID(clientID domain.ID) *Container {
	for _, ctr := range mgr.Containers {
		if ctr.GetID() == clientID {
			return ctr
		}
	}

	return nil
}

// FindContainer return the container associated to the client.
func (mgr *Manager) FindContainer(containerID domain.ID) *Container {
	if ctr, ok := mgr.Containers[containerID]; ok {
		return ctr
	}

	return nil
}
