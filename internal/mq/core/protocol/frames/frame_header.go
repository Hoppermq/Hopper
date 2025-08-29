package frames

import (
	"github.com/hoppermq/hopper/pkg/domain"
)

// Header represent the base frame header.
type Header struct {
	Size    uint16
	Type    domain.FrameType
	DOFF    domain.DOFF
	Channel uint8
}

// PayloadHeader represent the base payloadHeader.
type PayloadHeader struct {
	Size uint16
}

// GetFrameType return the type frame.
func (h *Header) GetFrameType() domain.FrameType {
	return h.Type
}

// Validate validate a frame.
func (h *Header) Validate() bool {
	return h.Type != 0 && h.DOFF != 0
}

// SetSize set the frame size.
func (h *Header) SetSize(s uint16) {
	h.Size = s
}

// GetSize return the current size.
func (h *Header) GetSize() uint16 {
	return h.Size
}

// GetDOFF return the frame doff.
func (h *Header) GetDOFF() domain.DOFF {
	return h.DOFF
}

// Sizer return the size.
func (ph *PayloadHeader) Sizer() uint16 {
	return 2
}
