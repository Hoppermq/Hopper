// Package domain provides the definitions and interfaces for the HopperMQ protocol frames.
package domain

// FrameType represents the type of frame in the HopperMQ protocol.
type FrameType uint8

// DOFF represents the Data Offset in the HopperMQ protocol.
type DOFF uint8

type Frame interface {
	GetType() FrameType
	GetHeader() HeaderFrame
	GetPayload() Payload
}

// HeaderFrame is the interface for all header frames in the HopperMQ protocol.
type HeaderFrame interface {
	Validate() bool
	GetFrameType() FrameType
	GetSize() uint16
	GetDOFF() DOFF
	SetSize(uint16)
}

type HeaderPayload interface {
	Sizer() uint16
}

// Payload is the interface for all payloads in the HopperMQ protocol.
type Payload interface {
	GetHeader() HeaderPayload
	GetData() []byte
	Sizer() uint16
}

// OpenFramePayload is the interface for open frame payloads in the HopperMQ protocol.
type OpenFramePayload interface {
	Payload
	GetSourceID() string
}

// MessageFramePayload is the interface for message frame payloads in the HopperMQ protocol.
type MessageFramePayload interface {
	Payload
}

const (
	// DOFF2 is the Data Offset for 2 bytes.
	DOFF2 DOFF = 2
	// DOFF3 is the Data Offset for 3 bytes.
	DOFF3 DOFF = 3
	// DOFF4 is the Data Offset for 4 bytes.
	DOFF4 DOFF = 4
)

const (
	// FrameTypeOpen is the frame type for open frames.
	FrameTypeOpen FrameType = 0x01

	// FrameTypeClose is the frame type for close frames.
	FrameTypeClose FrameType = 0x02

	// FrameTypeMessage is the frame type for message frames.
	FrameTypeMessage FrameType = 0x03

	// FrameTypeError is the frame type for error frames.
	FrameTypeError FrameType = 0xF0
)
