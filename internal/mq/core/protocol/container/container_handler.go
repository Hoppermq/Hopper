package container

import "github.com/hoppermq/hopper/pkg/domain"

func (ctr *Container) HandleFrame(f domain.Frame) error {
	switch ctr.State {
	case domain.ContainerOpenRcvd:
		println("received open rcvd")
	}
	return nil
}
