package frames

type OpenFrame struct {
	Frame
}

type OpenFramePayloadData struct {
	SourceID       string
	AssignedChanID string
}

func (op *OpenFramePayloadData) GetSourceID() string {
	return op.SourceID
}

func CreateOpenFramePayloadData(sourceID string, assignedChannel string) *OpenFramePayloadData {
	return &OpenFramePayloadData{
		SourceID:       sourceID,
		AssignedChanID: assignedChannel,
	}
}
