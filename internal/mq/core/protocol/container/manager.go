package container

import (
	"sync"

	"github.com/hoppermq/hopper/pkg/domain"
)

type ContainerManager struct {
	Container map[string]domain.Container

	mut sync.RWMutex
}

func NewContainerManager() *ContainerManager {
	return &ContainerManager{
		Container: make(map[string]domain.Container),
	}
}
