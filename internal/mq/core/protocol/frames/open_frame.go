package frames

import (
	"github.com/hoppermq/hopper/pkg/domain"
)

// OpenFramePayload represents the payload for open frames in the HopperMQ protocol.
type OpenFramePayload struct {
	BasePayload
	SourceID            domain.ID
	AssignedContainerID domain.ID
}

type OpenRcvdPayload struct {
	BasePayload
	SourceID domain.ID
}

// CreateOpenFramePayload creates a new OpenFramePayload instance.
func CreateOpenFramePayload(header domain.HeaderPayload, sourceID domain.ID, assignedContainerID domain.ID) *OpenFramePayload {
	return &OpenFramePayload{
		BasePayload: BasePayload{
			Header: header,
		},
		SourceID:            sourceID,
		AssignedContainerID: assignedContainerID,
	}
}

// GetSourceID returns the source ID from the open frame payload.
func (ofp *OpenFramePayload) GetSourceID() domain.ID {
	return ofp.SourceID
}

// GetAssignedContainerID returns the assigned container ID from the open frame payload.
func (ofp *OpenFramePayload) GetAssignedContainerID() domain.ID {
	return ofp.AssignedContainerID
}

// Sizer calculates the total size of the open frame payload.
func (ofp *OpenFramePayload) Sizer() uint16 {
	headerSize := uint16(0)
	if ofp.Header != nil {
		headerSize = ofp.Header.Sizer()
	}

	// Calculate size of IDs (assuming they serialize to known sizes)
	dataSize := uint16(len(ofp.SourceID) + len(ofp.AssignedContainerID))
	return headerSize + dataSize
}

func (o OpenRcvdPayload) Sizer() uint16 {
	headerSize := uint16(0)
	if o.Header != nil {
		headerSize = o.Header.Sizer()
	}

	dataSize := uint16(len(o.SourceID))
	return headerSize + dataSize
}

func (o OpenRcvdPayload) GetSourceID() domain.ID {
	return o.SourceID
}

