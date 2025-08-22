package container

import (
	"sync"

	"github.com/hoppermq/hopper/pkg/domain"
)

type ContainerManager struct {
	Container map[domain.ID]domain.Container

	mut sync.RWMutex
}

func NewContainerManager() *ContainerManager {
	return &ContainerManager{
		Container: make(map[domain.ID]domain.Container),
	}
}

func (ctnrManager *ContainerManager) CreateNewContainer(
	IDGenerator func() domain.ID,
	clientID domain.ID,
) domain.Container {
	println("HELLO ??")
	container := NewContainer(IDGenerator(), clientID)
	ctnrManager.Container[container.GetID()] = container

	return container
}
