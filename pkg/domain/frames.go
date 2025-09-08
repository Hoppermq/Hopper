// Package domain provides the definitions and interfaces for the HopperMQ protocol frames.
package domain

// FrameType represents the type of frame in the HopperMQ protocol.
type FrameType uint8

// DOFF represents the Data Offset in the HopperMQ protocol.
type DOFF uint8

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

	// FrameTypeOpenRcvd is the frame type for open received frames.
	FrameTypeOpenRcvd FrameType = 0x02

	// FrameTypeClose is the frame type for close frames.
	FrameTypeClose FrameType = 0x03

	// FrameTypeConnect is the frame type for connect frames.
	FrameTypeConnect FrameType = 0x04

	// FrameTypeSubscribe is the frame type for subscribe frames.
	FrameTypeSubscribe FrameType = 0x05

	// FrameTypeUnsubscribe is the frame type for unsubscribe frames.
	FrameTypeUnsubscribe FrameType = 0x06

	// FrameTypeAuth represent the frame type for authentification.
	FrameTypeAuth FrameType = 0x07

	// FrameTypeBegin represent the frame type for begin a connection.
	FrameTypeBegin FrameType = 0x08

	// FrameTypeStart represent the frame type for starting the message flow.
	FrameTypeStart FrameType = 0x0A

	// FrameTypeMessage represent the frame type for a message.
	FrameTypeMessage FrameType = 0x1F

	// FrameTypeError is the frame type for error frames.
	FrameTypeError FrameType = 0xF0
)

// Frame represent a frame used by the protocol.
type Frame interface {
	GetType() FrameType
	GetHeader() HeaderFrame
	GetPayload() Payload
	CanHandle(frameType FrameType) bool
}

// HeaderFrame is the interface for all header frames in the HopperMQ protocol.
type HeaderFrame interface {
	Validate() bool
	GetFrameType() FrameType
	GetSize() uint16
	GetDOFF() DOFF
	SetSize(uint16)
}

// HeaderPayload represent the domain interface of a frame payload header.
type HeaderPayload interface {
	Sizer() uint16
}

// Payload is the interface for all payloads in the HopperMQ protocol.
type Payload interface {
	GetHeader() HeaderPayload
	Sizer() uint16
}

// OpenFramePayload is the interface for open frame payloads in the HopperMQ protocol.
type OpenFramePayload interface {
	Payload
	GetSourceID() ID
	GetAssignedContainerID() ID
}

type OpenRcvdFramePayload interface {
	Payload
	GetSourceID() ID
}

// MessageFramePayload is the interface for message frame payloads in the HopperMQ protocol.
type MessageFramePayload interface {
	Payload
	GetTopic() string
	GetMessageID() ID
	GetSourceID() ID
	GetContent() []byte
	GetHeaders() map[string]string
}

// ConnectFramePayload is the interface for connect frame payloads in the HopperMQ protocol.
type ConnectFramePayload interface {
	Payload
	GetSourceID() ID
	GetClientVersion() string
	GetKeepAlive() uint16
}

// SubscribeFramePayload is the interface for subscribe frame payloads in the HopperMQ protocol.
type SubscribeFramePayload interface {
	Payload
	GetTopic() string
	GetQoS() uint8
	GetRoutingKey() string
}

// UnsubscribeFramePayload is the interface for unsubscribe frame payloads in the HopperMQ protocol.
type UnsubscribeFramePayload interface {
	Payload
	GetTopic() string
}

// CloseFramePayload is the interface for close frame payloads in the HopperMQ protocol.
type CloseFramePayload interface {
	Payload
	GetReason() string
	GetCode() uint16
}

// ErrorFramePayload is the interface for error frame payloads in the HopperMQ protocol.
type ErrorFramePayload interface {
	Payload
	GetErrorCode() uint16
	GetErrorMessage() string
	GetDetails() map[string]string
}

// BeginFramePayload is the interface for begin frame payloads in the HopperMQ protocol.
type BeginFramePayload interface {
	Payload
	GetSourceID() ID
	GetContainerID() ID
	GetRemoteChannel() uint16
	GetNextOutgoingID() uint32
	GetIncomingWindow() uint32
	GetOutgoingWindow() uint32
}
