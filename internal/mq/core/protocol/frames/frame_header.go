package frames

import (
	"github.com/hoppermq/hopper/pkg/domain"
)

type Header struct {
	Size    uint16
	Type    domain.FrameType
	DOFF    domain.DOFF
	Channel uint8
}

type PayloadHeader struct {
	Size uint16
}

func (h *Header) GetFrameType() domain.FrameType {
	return h.Type
}

func (h *Header) Validate() bool {
	panic("implement me")
}

func (h *Header) SetSize(s uint16) {
	h.Size = s
}

func (ph *PayloadHeader) Sizer() uint16 {
	return 2
}
