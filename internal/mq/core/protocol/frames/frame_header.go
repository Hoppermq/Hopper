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
	return h.Type != 0 && h.DOFF != 0
}

func (h *Header) SetSize(s uint16) {
	h.Size = s
}

func (h *Header) GetSize() uint16 {
	return h.Size
}

func (h *Header) GetDOFF() domain.DOFF {
	return h.DOFF
}

func (ph *PayloadHeader) Sizer() uint16 {
	return 2
}
