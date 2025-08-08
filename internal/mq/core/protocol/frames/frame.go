// Package frames is the protocol frames business logic for HopperMQ.
package frames

import "github.com/hoppermq/hopper/pkg/domain"

// ExtendedFrameHeader is an interface for extended frame headers in the HopperMQ protocol.
type ExtendedFrameHeader interface {
	domain.Serializable
}

type Header struct {
	Size      uint16
	DOFF      domain.DOFF
	FrameType domain.FrameType
}

type PayloadHeader struct {
	Size uint16
}

type Payload struct {
	Header domain.HeaderPayload
	Data   []byte
}

// Frame represents a protocol frame in the HopperMQ system.
type Frame struct {
	domain.Serializable
	Header         domain.HeaderFrame
	ExtendedHeader ExtendedFrameHeader
	Payload        domain.Payload
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
) *Frame {
	if err := validateFrame(header, payload); err != nil {
		return nil
	}

	return &Frame{
		Header:         header,
		ExtendedHeader: extendedHeader,
		Payload:        payload,
	}
}

func (f *Frame) Serialize() ([]byte, error) {
	// Implement serialization logic here
	return nil, nil
}

func (f *Frame) Deserialize(data []byte) (domain.Serializable, error) {
	// Implement deserialization logic here
	return nil, nil
}

func (p *Payload) Serialize() ([]byte, error) {
	// Implement serialization logic for Payload
	return nil, nil
}

func (p *Payload) Deserialize(data []byte) (domain.Serializable, error) {
	// Implement deserialization logic for Payload
	return nil, nil
}

func (ph *PayloadHeader) Serialize() ([]byte, error) {
	// Implement serialization logic for PayloadHeader
	return nil, nil
}

func (ph *PayloadHeader) Deserialize(data []byte) (domain.Serializable, error) {
	// Implement deserialization logic for PayloadHeader
	return nil, nil
}
