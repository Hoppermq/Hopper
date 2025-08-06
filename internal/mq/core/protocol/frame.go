package protocol

type FrameHeader struct {
	Size uint32
	Type string
}

type Frame struct {
	Header *FrameHeader

	body struct{}
}
