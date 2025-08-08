package frames

import "github.com/hoppermq/hopper/pkg/domain"

type HeaderFrame struct {
	Size    uint32
	Type    domain.FrameType
	DOFF    domain.DOFF
	Channel uint8
}

func (h *HeaderFrame) Serialize() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (h *HeaderFrame) Deserialize(data []byte) (domain.Serializable, error) {
	//TODO implement me
	panic("implement me")
}

func (h *HeaderFrame) GetFrameType() domain.FrameType {
	return h.Type
}

func (h *HeaderFrame) Validate() bool {
	//TODO implement me
	panic("implement me")
}
