package container

import (
	"sync"

	"github.com/hoppermq/hopper/pkg/domain"
)

type ContainerRegistry struct {
	mu sync.RWMutex

	data map[string]map[domain.ID]struct{}
}
type ContainerManager struct {
	Registry ContainerRegistry

	mut sync.RWMutex
}

func NewContainerRegistry() *ContainerRegistry {
	return &ContainerRegistry{
		data: make(map[string]map[domain.ID]struct{}),
	}
}

func NewContainerManager() *ContainerManager {
	return &ContainerManager{
		Registry: *NewContainerRegistry(),
	}
}

func (ctnrManager *ContainerManager) CreateNewContainer(
	IDGenerator func() domain.ID,
	clientID domain.ID,
) domain.Container {
	container := NewContainer(IDGenerator(), clientID)

	return container
}

func (rContainer *ContainerRegistry) Register(topic string, id domain.ID) {
	rContainer.mu.Lock()
	defer rContainer.mu.Unlock()

	if rContainer.data[topic] == nil {
		rContainer.data[topic] = make(map[domain.ID]struct{})
	}

	rContainer.data[topic][id] = struct{}{}
}

func (rContainer *ContainerRegistry) Unregister(topic string, id domain.ID) {
	rContainer.mu.Lock()
	defer rContainer.mu.Unlock()

	if set, ok := rContainer.data[topic]; ok {
		delete(set, id)
		if len(set) == 0 {
			delete(rContainer.data, topic)
		}
	}
}

func (ctnrManager *ContainerManager) RegisterContainerToTopic(
	topic string,
	containerID domain.ID,
) {
	ctnrManager.Registry.Register(topic, containerID)
}

func (ctnrManager *ContainerManager) RemoveContainerFromTopic(
	topic string,
	containerID domain.ID,
) {
	ctnrManager.Registry.Unregister(topic, containerID)
}
