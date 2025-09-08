// Package frames is the protocol frames business logic for HopperMQ.
package frames

import (
	"github.com/hoppermq/hopper/internal/common"
	"github.com/hoppermq/hopper/pkg/domain"
)

// ExtendedFrameHeader is an interface for extended frame headers in the HopperMQ protocol.
type ExtendedFrameHeader interface {
}

// BasePayload provides common payload functionality.
type BasePayload struct {
	Header domain.HeaderPayload
}

// Frame represents a protocol frame in the HPMQ Protocol.
type Frame struct {
	Header         domain.HeaderFrame
	ExtendedHeader ExtendedFrameHeader
	Payload        domain.Payload
}

// GetHeader return the frame header.
func (f Frame) GetHeader() domain.HeaderFrame {
	return f.Header
}

// GetPayload return the frame payload.
func (f Frame) GetPayload() domain.Payload {
	return f.Payload
}

// GetType return the frame type.
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
	case domain.FrameTypeConnect:
		if _, ok := payload.(domain.ConnectFramePayload); !ok {
			return domain.ErrInvalidPayload
		}
	case domain.FrameTypeSubscribe:
		if _, ok := payload.(domain.SubscribeFramePayload); !ok {
			return domain.ErrInvalidPayload
		}
	case domain.FrameTypeUnsubscribe:
		if _, ok := payload.(domain.UnsubscribeFramePayload); !ok {
			return domain.ErrInvalidPayload
		}
	case domain.FrameTypeClose:
		if _, ok := payload.(domain.CloseFramePayload); !ok {
			return domain.ErrInvalidPayload
		}
	case domain.FrameTypeError:
		if _, ok := payload.(domain.ErrorFramePayload); !ok {
			return domain.ErrInvalidPayload
		}
	case domain.FrameTypeBegin:
		if _, ok := payload.(domain.BeginFramePayload); !ok {
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
) (*Frame, error) {
	if err := validateFrame(header, payload); err != nil {
		return nil, err
	}

	header.SetSize(calculatePayloadSize(payload))

	return &Frame{
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

// GetHeader return the frame payload header.
func (bp *BasePayload) GetHeader() domain.HeaderPayload {
	return bp.Header
}
