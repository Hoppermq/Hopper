// Package frames is the protocol frames business logic for HopperMQ.
package frames

import (
	"github.com/hoppermq/hopper/internal/common"
	"github.com/hoppermq/hopper/pkg/domain"
)

// ExtendedFrameHeader is an interface for extended frame headers in the HopperMQ protocol.
type ExtendedFrameHeader interface {
}

// Payload represent the content of the frame in the HPMQ Protocol.
type Payload struct {
	Header domain.HeaderPayload
	Data   []byte
}

// Frame represents a protocol frame in the HPMQ Protocol.
type Frame struct {
	Header         domain.HeaderFrame
	ExtendedHeader ExtendedFrameHeader
	Payload        domain.Payload
}

func (f Frame) GetHeader() domain.HeaderFrame {
	return f.Header
}

func (f Frame) GetPayload() domain.Payload {
	return f.Payload
}

func (f Frame) GetType() domain.FrameType {
	return f.Header.GetFrameType()
}

func validateFrame(header domain.HeaderFrame, payload domain.Payload) error {
	if header == nil {
		return domain.ErrInvalidHeader
	}
	frameType := header.GetFrameType()
	switch frameType {
	case domain.FrameTypeOpen:
		if _, ok := payload.(domain.OpenFramePayload); !ok {
			return domain.ErrInvalidPayload
		}
	case domain.FrameTypeMessage:
		if _, ok := payload.(domain.MessageFramePayload); !ok {
			return domain.ErrInvalidPayload
		}
	case domain.FrameTypeClose:
		if _, ok := payload.(domain.OpenFramePayload); !ok {
			return domain.ErrInvalidPayload
		}
	default:
		return nil
	}
	return nil
}

// CreateFrame creates a new Frame with the given header, extended header, and payload.
func CreateFrame(
	header domain.HeaderFrame,
	extendedHeader ExtendedFrameHeader,
	payload domain.Payload,
) (Frame, error) {
	/* if err := validateFrame(header, payload); err != nil {
	return Frame{}, err
	}*/

	header.SetSize(calculatePayloadSize(payload))

	return Frame{
		Header:         header,
		ExtendedHeader: extendedHeader,
		Payload:        payload,
	}, nil
}

func calculatePayloadSize(payload domain.Payload) uint16 {
	if sizer, ok := payload.(interface{ Sizer() uint16 }); ok {
		return sizer.Sizer()
	}

	if data, err := common.Serialize(payload); err == nil {
		return uint16(len(data))
	}

	return 0
}

func (p *Payload) Sizer() uint16 {
	headerSize := uint16(0)
	if p.Header != nil {
		if sizer, ok := p.Header.(interface{ Sizer() uint16 }); ok {
			headerSize = sizer.Sizer()
		} else {
			if data, err := common.Serialize(p.Header); err != nil {
				headerSize = uint16(len(data))
			}
		}
	}

	dataSize := uint16(len(p.Data))
	return headerSize + dataSize
}

func (p *Payload) GetHeader() domain.HeaderPayload {
	return p.Header
}

func (p *Payload) GetData() []byte {
	return p.Data
}
