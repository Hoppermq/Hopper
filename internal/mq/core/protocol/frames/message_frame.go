package frames

import (
	"github.com/hoppermq/hopper/pkg/domain"
)

// MessageFramePayload represents the payload for message frames in the HopperMQ protocol.
type MessageFramePayload struct {
	BasePayload
	Topic     string
	SourceID  domain.ID
	MessageID domain.ID
	Content   []byte
	Headers   map[string]string
}

func (mfp *MessageFramePayload) GetSourceID() domain.ID {
	return mfp.SourceID
}

// GetTopic returns the topic from the message frame payload.
func (mfp *MessageFramePayload) GetTopic() string {
	return mfp.Topic
}

// GetMessageID returns the message ID from the message frame payload.
func (mfp *MessageFramePayload) GetMessageID() domain.ID {
	return mfp.MessageID
}

// GetContent returns the content from the message frame payload.
func (mfp *MessageFramePayload) GetContent() []byte {
	return mfp.Content
}

// GetHeaders returns the headers from the message frame payload.
func (mfp *MessageFramePayload) GetHeaders() map[string]string {
	return mfp.Headers
}

// Sizer calculates the total size of the message frame payload.
func (mfp *MessageFramePayload) Sizer() uint16 {
	headerSize := uint16(0)
	if mfp.Header != nil {
		headerSize = mfp.Header.Sizer()
	}

	dataSize := uint16(len(mfp.Topic) + len(mfp.MessageID) + len(mfp.Content))

	for k, v := range mfp.Headers {
		dataSize += uint16(len(k) + len(v))
	}

	return headerSize + dataSize
}

// CreateMessageFramePayload creates a new MessageFramePayload instance.
func CreateMessageFramePayload(
	header domain.HeaderPayload,
	topic string,
	messageID domain.ID,
	content []byte,
	headers map[string]string,
) *MessageFramePayload {
	if headers == nil {
		headers = make(map[string]string)
	}

	return &MessageFramePayload{
		BasePayload: BasePayload{
			Header: header,
		},
		Topic:     topic,
		MessageID: messageID,
		Content:   content,
		Headers:   headers,
	}
}
