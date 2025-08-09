package frames

import (
	"github.com/hoppermq/hopper/pkg/domain"
)

type HeaderFrame struct {
	Size    uint16
	Type    domain.FrameType
	DOFF    domain.DOFF
	Channel uint8
}

func (h *HeaderFrame) GetFrameType() domain.FrameType {
	return h.Type
}

func (h *HeaderFrame) Validate() bool {
	panic("implement me")
}

func (h *HeaderFrame) SetSize(s uint16) {
	h.Size = s
}
