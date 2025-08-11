package frames

import (
	"github.com/hoppermq/hopper/pkg/domain"
)

func CreateOpenFrame(doff domain.DOFF) (domain.Frame, error) {
	headerFrame := Header{
		Size: 0,
		DOFF: doff,
		Type: domain.FrameTypeOpen,
	}

	payload := Payload{
		Header: &PayloadHeader{
			Size: 0,
		},
		Data: []byte("OPEN FRAME"),
	}

	return CreateFrame(&headerFrame, nil, &payload)
}

type OpenFrame struct {
	Frame

	SourceID string
}

func (op *OpenFrame) GetSourceID() string {
	return op.SourceID
}
