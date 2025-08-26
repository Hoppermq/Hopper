package frames

import (
	"github.com/hoppermq/hopper/pkg/domain"
)

func CreateOpenFrame(doff domain.DOFF, sourceID domain.ID, assignedContainerID domain.ID) (domain.Frame, error) {
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

func CreateMessageFrame(
	doff domain.DOFF,
	topic string,
	messageID domain.ID,
	content []byte,
	headers map[string]string,
) (domain.Frame, error) {
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
