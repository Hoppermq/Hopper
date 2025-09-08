package frames

import "github.com/hoppermq/hopper/pkg/domain"

// BeginFramePayload represent the Begin Frame Payload.
type BeginFramePayload struct {
	BasePayload
	SourceID       domain.ID
	ContainerID    domain.ID
	RemoteChannel  uint16
	NextOutgoingID uint32
	IncomingWindow uint32
	OutgoingWindow uint32
}

// Sizer return the payload size.
func (f *BeginFramePayload) Sizer() uint16 {
	headerSize := uint16(0)
	if f.Header != nil {
		headerSize = f.Header.Sizer()
	}

	dataSize := uint16(len(f.SourceID) + len(f.ContainerID) + 2 + 4 + 4 + 4)

	return headerSize + dataSize
}

// GetSourceID return the source ID.
func (f *BeginFramePayload) GetSourceID() domain.ID {
	return f.SourceID
}

// GetContainerID return the container ID.
func (f *BeginFramePayload) GetContainerID() domain.ID {
	return f.ContainerID
}

// GetRemoteChannel return the remote channel number.
func (f *BeginFramePayload) GetRemoteChannel() uint16 {
	return f.RemoteChannel
}

// GetNextOutgoingID return the next outgoing ID.
func (f *BeginFramePayload) GetNextOutgoingID() uint32 {
	return f.NextOutgoingID
}

// GetIncomingWindow return the incoming window size.
func (f *BeginFramePayload) GetIncomingWindow() uint32 {
	return f.IncomingWindow
}

// GetOutgoingWindow return the outgoing window size.
func (f *BeginFramePayload) GetOutgoingWindow() uint32 {
	return f.OutgoingWindow
}

// CreateBeginFramePayload creates a new BeginFramePayload instance.
func CreateBeginFramePayload(
	header *PayloadHeader,
	sourceID domain.ID,
	containerID domain.ID,
	remoteChannel uint16,
	nextOutgoingID uint32,
	incomingWindow uint32,
	outgoingWindow uint32,
) *BeginFramePayload {
	return &BeginFramePayload{
		BasePayload: BasePayload{
			Header: header,
		},
		SourceID:       sourceID,
		ContainerID:    containerID,
		RemoteChannel:  remoteChannel,
		NextOutgoingID: nextOutgoingID,
		IncomingWindow: incomingWindow,
		OutgoingWindow: outgoingWindow,
	}
}
