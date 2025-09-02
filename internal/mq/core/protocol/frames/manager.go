package frames

import "github.com/hoppermq/hopper/pkg/domain"

// FrameManager will classify the frame.
type FrameManager struct {}

// IsControlFrame return if frame is in range of control frame.
func (fm *FrameManager) IsControlFrame(ft domain.FrameType) bool {
	return ft >= 0x01 && ft <= 0x0F
}

// IsMessageFrame return if frame type is in range of message frame.
func (fm *FrameManager) IsMessageFrame(ft domain.FrameType) bool {
	return ft >= 0x10 && ft <= 0x1F
}

// IsErrorFrame return if frame type is in range of error frame.
func (fm *FrameManager) IsErrorFrame(ft domain.FrameType) bool {
	return ft >= 0xf0
}
