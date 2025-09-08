package frames

import "github.com/hoppermq/hopper/pkg/domain"

// ConnectFramePayload represent the Connect Frame Payload.
type ConnectFramePayload struct {
	BasePayload
	SourceID      domain.ID
	clientVersion string
	keepAlive     uint16
}

// Sizer return the payload size.
func (f *ConnectFramePayload) Sizer() uint16 {
	headerSize := uint16(0)
	if f.Header != nil {
		headerSize = f.Header.Sizer()
	}

	dataSize := uint16(4 + len(f.SourceID) + len(f.clientVersion)) // missing len of uint16.

	return headerSize + dataSize
}

// GetSourceID return the source ID.
func (f *ConnectFramePayload) GetSourceID() domain.ID {
	return f.SourceID
}

// GetClientVersion return the client version.
func (f *ConnectFramePayload) GetClientVersion() string {
	return f.clientVersion
}

// GetKeepAlive return the channel lt set by the client.
func (f *ConnectFramePayload) GetKeepAlive() uint16 {
	return f.keepAlive
}
