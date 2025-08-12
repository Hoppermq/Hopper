package serializer

import (
	"bytes"
	"encoding/binary"
	"sync"

	"github.com/hoppermq/hopper/internal/mq/core/protocol/frames"
	"github.com/hoppermq/hopper/pkg/domain"
)

type Serializer struct {
	bufferPool domain.Pool[*bytes.Buffer] // since it's a global Serializer we will keep it as any atm.
	mu         sync.RWMutex
}

func (ps *Serializer) writeUint16(b *bytes.Buffer, u16 uint16) error {
	return binary.Write(b, binary.BigEndian, u16)
}

func (ps *Serializer) writeUint32(b *bytes.Buffer, u32 uint32) error {
	return binary.Write(b, binary.BigEndian, u32)
}

func (ps *Serializer) writeByteArray(b *bytes.Buffer, d []byte) error {
	if err := ps.writeUint32(b, uint32(len(d))); err != nil {
		return err
	}
	_, err := b.Write(d)
	return err
}

func (ps *Serializer) writeString(b *bytes.Buffer, str string) error {
	if err := ps.writeUint32(b, uint32(len(str))); err != nil {
		return err
	}

	_, err := b.Write([]byte(str))
	return err
}

func (ps *Serializer) writeFrameHeader(buff *bytes.Buffer, fh domain.HeaderFrame) error {
	if err := ps.writeUint16(buff, fh.GetSize()); err != nil {
		return err
	}
	if err := ps.writeUint16(buff, uint16(fh.GetDOFF())); err != nil {
		return err
	}
	return ps.writeUint16(buff, uint16(fh.GetFrameType()))
}

func (ps *Serializer) writePayloadHeader(buff *bytes.Buffer, ph domain.HeaderPayload) error {
	return ps.writeUint16(buff, ph.Sizer())
}

func (ps *Serializer) writePayload(buff *bytes.Buffer, fp domain.Payload) error {
	if err := ps.writePayloadHeader(buff, fp.GetHeader()); err != nil {
		return err
	}
	return ps.writeByteArray(buff, fp.GetData())
}

func (ps *Serializer) SerializeFrame(
	frame domain.Frame,
) ([]byte, error) {
	buff := ps.bufferPool.Get()
	defer ps.bufferPool.Put(buff)
	buff.Reset()

	if err := ps.writeFrameHeader(buff, frame.GetHeader()); err != nil {
		return nil, err
	}

	if err := ps.writePayload(buff, frame.GetPayload()); err != nil {
		return nil, err
	}

	res := make([]byte, buff.Len())
	copy(res, buff.Bytes())

	return res, nil
}

func (ps *Serializer) SerializeOpenFramePayloadData(data *frames.OpenFramePayloadData) ([]byte, error) {
	buff := ps.bufferPool.Get()
	defer ps.bufferPool.Put(buff)
	buff.Reset()

	if err := ps.writeString(buff, data.AssignedChanID); err != nil {
		return nil, err
	}

	if err := ps.writeString(buff, data.SourceID); err != nil {
		return nil, err
	}

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

func NewSerializer(pool domain.Pool[*bytes.Buffer]) *Serializer {
	return &Serializer{
		bufferPool: pool,
	}
}
