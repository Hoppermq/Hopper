package container

import "github.com/hoppermq/hopper/pkg/domain"

func (ctnr *Container) HandleFrame(f domain.Frame) error {
	switch ctnr.State {
	case domain.OpenRcvd:
		
	}
	return  nil
}
