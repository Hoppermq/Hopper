// Package frames is the protocol frames business logic for HopperMQ.
package frames

import "github.com/hoppermq/hopper/pkg/domain"

// ExtendedFrameHeader is an interface for extended frame headers in the HopperMQ protocol.
type ExtendedFrameHeader interface {
	domain.Serializable
}

// Frame represents a protocol frame in the HopperMQ system.
type Frame struct {
	domain.Serializable
	Header         *domain.HeaderFrame
	ExtendedHeader ExtendedFrameHeader
	Payload        *domain.Payload
}

// CreateFrame creates a new Frame with the given header, extended header, and payload.
func CreateFrame(header *domain.HeaderFrame, extendedHeader ExtendedFrameHeader, payload *domain.Payload) *Frame {
	// We should probably validate the header and payload here
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
