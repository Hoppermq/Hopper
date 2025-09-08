package container

import "github.com/hoppermq/hopper/pkg/domain"

// UpdateContainerState will update the container state to a new state.
func (manager *Manager) UpdateContainerState(containerID domain.ID, newState domain.ContainerState) {
	manager.Containers[containerID].SetState(newState)
}
