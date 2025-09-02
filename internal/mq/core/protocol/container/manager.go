package container

import (
	"sync"

	"github.com/hoppermq/hopper/pkg/domain"
)

// Registry represent the local registry for the orchestrator.
type Registry struct {
	mu sync.RWMutex

	data map[string]map[domain.ID]struct{}
}

// Manager represent the container orchestrator.
type Manager struct {
	Registry   *Registry
	Containers map[domain.ID]*Container

	mut sync.RWMutex
}

// NewContainerRegistry return a new registry.
func NewContainerRegistry() *Registry {
	return &Registry{
		data: make(map[string]map[domain.ID]struct{}),
	}
}

// NewContainerManager return a new instance of the container orchestrator.
func NewContainerManager() *Manager {
	return &Manager{
		Registry:   NewContainerRegistry(),
		Containers: make(map[domain.ID]*Container),
	}
}

// CreateNewContainer create a new container.
func (ctnrManager *Manager) CreateNewContainer(
	idGenerator func() domain.ID,
	clientID domain.ID,
) *Container {
	container := NewContainer(idGenerator(), clientID)
	ctnrManager.Containers[container.ID] = container

	return container
}

// Register register attach a topic to a containerID.
func (rContainer *Registry) Register(topic string, containerID domain.ID) {
	rContainer.mu.Lock()
	defer rContainer.mu.Unlock()

	if rContainer.data[topic] == nil {
		rContainer.data[topic] = make(map[domain.ID]struct{})
	}

	rContainer.data[topic][containerID] = struct{}{}
}

// Unregister remove a topic from a container.
func (rContainer *Registry) Unregister(topic string, id domain.ID) {
	rContainer.mu.Lock()
	defer rContainer.mu.Unlock()

	if set, ok := rContainer.data[topic]; ok {
		delete(set, id)
		if len(set) == 0 {
			delete(rContainer.data, topic)
		}
	}
}

// RegisterContainerToTopic set a container to the registry attached to a topic.
func (ctnrManager *Manager) RegisterContainerToTopic(
	topic string,
	containerID domain.ID,
) {
	ctnrManager.Registry.Register(topic, containerID)
}

// RemoveContainerFromTopic remove the container from the registry to it's given topic.
func (ctnrManager *Manager) RemoveContainerFromTopic(
	topic string,
	containerID domain.ID,
) {
	ctnrManager.Registry.Unregister(topic, containerID)
}

