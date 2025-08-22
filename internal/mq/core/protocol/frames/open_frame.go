package frames

import "github.com/hoppermq/hopper/pkg/domain"

type OpenFrame struct {
	Frame
}

type OpenFramePayloadData struct {
	SourceID            domain.ID
	AssignedContainerID domain.ID
}

func (op *OpenFramePayloadData) GetSourceID() domain.ID {
	return op.SourceID
}

func CreateOpenFramePayloadData(sourceID domain.ID, assignedContainerID domain.ID) *OpenFramePayloadData {
	return &OpenFramePayloadData{
		SourceID:            sourceID,
		AssignedContainerID: assignedContainerID,
	}
}
