package frames

import (
	"github.com/hoppermq/hopper/pkg/domain"
)

func CreateOpenFrame(doff domain.DOFF, serializedData []byte) (domain.Frame, error) {
	headerFrame := Header{
		Size: 0,
		DOFF: doff,
		Type: domain.FrameTypeOpen,
	}

	payload := Payload{
		Header: &PayloadHeader{
			Size: 0,
		},
		Data: serializedData,
	}

	return CreateFrame(&headerFrame, nil, &payload)
}
