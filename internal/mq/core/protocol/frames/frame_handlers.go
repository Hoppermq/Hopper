package frames

import "github.com/hoppermq/hopper/pkg/domain"

func CreateOpenFrame(doff domain.DOFF) *Frame {
	return CreateFrame(
		&HeaderFrame{
			Size: 0, // Size will be set later
			DOFF: doff,
			Type: domain.FrameTypeOpen,
		},
		nil,
		&Payload{
			Header: &PayloadHeader{
				Size: 0, // Size will be set later
			},
		},
	)
}
