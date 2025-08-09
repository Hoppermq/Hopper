package frames

import (
	"github.com/hoppermq/hopper/pkg/domain"
)

func CreateOpenFrame(doff domain.DOFF) (Frame, error) {
	headerFrame := HeaderFrame{
		Size: 0,
		DOFF: doff,
		Type: domain.FrameTypeOpen,
	}

	payload := Payload{
		Header: &PayloadHeader{
			Size: 0,
		},
	}

	return CreateFrame(&headerFrame, nil, &payload)
}
