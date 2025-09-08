package frames

import (
	"github.com/hoppermq/hopper/pkg/domain"
)

// CreateOpenFrame create a new open frame.
func CreateOpenFrame(
	doff domain.DOFF,
	sourceID domain.ID,
	assignedContainerID domain.ID,
) (*Frame, error) {
	headerFrame := Header{
		Size: 0,
		DOFF: doff,
		Type: domain.FrameTypeOpen,
	}

	payloadHeader := &PayloadHeader{
		Size: 0,
	}

	payload := CreateOpenFramePayload(payloadHeader, sourceID, assignedContainerID)

	return CreateFrame(&headerFrame, nil, payload)
}

// CreateMessageFrame create a new MessageFrame.
func CreateMessageFrame(
	doff domain.DOFF,
	topic string,
	messageID domain.ID,
	content []byte,
	headers map[string]string,
) (*Frame, error) {
	headerFrame := Header{
		Size: 0,
		DOFF: doff,
		Type: domain.FrameTypeMessage,
	}

	payloadHeader := &PayloadHeader{
		Size: 0,
	}

	payload := CreateMessageFramePayload(payloadHeader, topic, messageID, content, headers)

	return CreateFrame(&headerFrame, nil, payload)
}

// CreateBeginFrame create a new begin frame.
func CreateBeginFrame(
	doff domain.DOFF,
	sourceID domain.ID,
	containerID domain.ID,
	remoteChannel uint16,
	nextOutgoingID uint32,
	incomingWindow uint32,
	outgoingWindow uint32,
) (*Frame, error) {
	headerFrame := Header{
		Size: 0,
		DOFF: doff,
		Type: domain.FrameTypeBegin,
	}

	payloadHeader := &PayloadHeader{
		Size: 0,
	}

	payload := CreateBeginFramePayload(
		payloadHeader,
		sourceID,
		containerID,
		remoteChannel,
		nextOutgoingID,
		incomingWindow,
		outgoingWindow,
	)

	return CreateFrame(&headerFrame, nil, payload)
}

// CanHandle return if frame match the frame type ?.
func (f *Frame) CanHandle(frameType domain.FrameType) bool {
	return f.GetType() == frameType
}
