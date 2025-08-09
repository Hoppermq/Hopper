package serializer

import (
	"bytes"
	"encoding/binary"
	"sync"

	"github.com/hoppermq/hopper/pkg/domain"
)

type Serializer struct {
	bufferPool domain.Pool[any] // since it's a global Serializer we will keep it as any atm.
	mu         sync.RWMutex
}

func (ps *Serializer) readUint16()        {}
func (ps *Serializer) readUint32()        {}
func (ps *Serializer) readFrameHeader()   {}
func (ps *Serializer) readFramePayload()  {}
func (ps *Serializer) readPayloadHeader() {}
func (ps *Serializer) readPayload()       {}

func (ps *Serializer) SerializeFrame(
	frame domain.Frame,
) ([]byte, error) {
	buff := ps.bufferPool.Get().(*bytes.Buffer)
	defer ps.bufferPool.Put(buff)
	buff.Reset()

	binary.Write(buff, binary.BigEndian, frame)
	res := make([]byte, buff.Len())
	copy(res, buff.Bytes())

	return res, nil
}

func (ps *Serializer) DeserializeFrame(d []byte) (domain.Frame, error) {
	r := bytes.NewReader(d)
	var f domain.Frame

	binary.Read(r, binary.BigEndian, &f)

	return nil, nil
}
