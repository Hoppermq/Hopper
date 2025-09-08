package serializer

import (
	"bytes"
	"encoding/binary"
	"sync"

	"github.com/hoppermq/hopper/internal/mq/core/protocol/frames"
	"github.com/hoppermq/hopper/pkg/domain"
)

// Serializer represent the protocol serializer.
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

	_, err := b.WriteString(str)
	return err
}

func (ps *Serializer) writeID(b *bytes.Buffer, ID domain.ID) error {
	if err := ps.writeUint32(b, uint32(len(ID))); err != nil {
		return err
	}

	_, err := b.WriteString(string(ID))
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

func (ps *Serializer) writePayload(buff *bytes.Buffer, frame domain.Frame) error {
	if err := ps.writePayloadHeader(buff, frame.GetPayload().GetHeader()); err != nil {
		return err
	}

	switch frame.GetType() {
	case domain.FrameTypeOpen:
		if openPayload, ok := frame.GetPayload().(domain.OpenFramePayload); ok {
			return ps.writeOpenPayload(buff, openPayload)
		}
	case domain.FrameTypeMessage:
		if msgPayload, ok := frame.GetPayload().(domain.MessageFramePayload); ok {
			return ps.writeMessagePayload(buff, msgPayload)
		}
	case domain.FrameTypeConnect:
		return domain.ErrUnsupportedFrameType
	default:
		return domain.ErrUnsupportedFrameType
	}
	return domain.ErrInvalidPayload
}

func (ps *Serializer) writeOpenPayload(buff *bytes.Buffer, payload domain.OpenFramePayload) error {
	if err := ps.writeID(buff, payload.GetSourceID()); err != nil {
		return err
	}
	return ps.writeID(buff, payload.GetAssignedContainerID())
}

func (ps *Serializer) writeMessagePayload(buff *bytes.Buffer, payload domain.MessageFramePayload) error {
	if err := ps.writeString(buff, payload.GetTopic()); err != nil {
		return err
	}
	if err := ps.writeID(buff, payload.GetMessageID()); err != nil {
		return err
	}
	if err := ps.writeByteArray(buff, payload.GetContent()); err != nil {
		return err
	}

	if err := ps.writeUint32(buff, uint32(len(payload.GetHeaders()))); err != nil {
		return err
	}
	for k, v := range payload.GetHeaders() {
		if err := ps.writeString(buff, k); err != nil {
			return err
		}
		if err := ps.writeString(buff, v); err != nil {
			return err
		}
	}
	return nil
}

func (ps *Serializer) writeConnectPayload(buff *bytes.Buffer, payload domain.ConnectFramePayload) error {
	if err := ps.writeID(buff, payload.GetSourceID()); err != nil {
		return err
	}
	if err := ps.writeString(buff, payload.GetClientVersion()); err != nil {
		return err
	}
	return ps.writeUint16(buff, payload.GetKeepAlive())
}

// SerializeFrame serialize the given frame.
func (ps *Serializer) SerializeFrame(frame domain.Frame) ([]byte, error) {
	buff := ps.bufferPool.Get()
	defer ps.bufferPool.Put(buff)
	buff.Reset()

	if err := ps.writeFrameHeader(buff, frame.GetHeader()); err != nil {
		return nil, err
	}

	if err := ps.writePayload(buff, frame); err != nil {
		return nil, err
	}

	res := make([]byte, buff.Len())
	copy(res, buff.Bytes())

	return res, nil
}

// DeserializeFrame deserialize the given bytes to frame.
func (ps *Serializer) DeserializeFrame(d []byte) (domain.Frame, error) {
	r := bytes.NewReader(d)

	var size, doff, frameType uint16
	if err := binary.Read(r, binary.BigEndian, &size); err != nil {
		return nil, err
	}
	if err := binary.Read(r, binary.BigEndian, &doff); err != nil {
		return nil, err
	}
	if err := binary.Read(r, binary.BigEndian, &frameType); err != nil {
		return nil, err
	}

	header := &frames.Header{
		Size: size,
		DOFF: domain.DOFF(doff),
		Type: domain.FrameType(frameType),
	}

	var payloadHeaderSize uint16
	if err := binary.Read(r, binary.BigEndian, &payloadHeaderSize); err != nil {
		return nil, err
	}

	payloadHeader := &frames.PayloadHeader{Size: payloadHeaderSize}

	var payload domain.Payload
	var err error

	switch domain.FrameType(frameType) {
	case domain.FrameTypeOpen:
		payload, err = ps.deserializeOpenPayload(r, payloadHeader)
	case domain.FrameTypeMessage:
		payload, err = ps.deserializeMessagePayload(r, payloadHeader)
	case domain.FrameTypeConnect:
		return nil, domain.ErrUnsupportedFrameType
	default:
		return nil, domain.ErrUnsupportedFrameType
	}

	if err != nil {
		return nil, err
	}

	return frames.CreateFrame(header, nil, payload)
}

func (ps *Serializer) deserializeOpenPayload(r *bytes.Reader, header domain.HeaderPayload) (*frames.OpenFramePayload, error) {
	sourceID, err := ps.readID(r)
	if err != nil {
		return nil, err
	}

	assignedContainerID, err := ps.readID(r)
	if err != nil {
		return nil, err
	}

	return frames.CreateOpenFramePayload(header, sourceID, assignedContainerID), nil
}

func (ps *Serializer) deserializeMessagePayload(r *bytes.Reader, header domain.HeaderPayload) (*frames.MessageFramePayload, error) {
	topic, err := ps.readString(r)
	if err != nil {
		return nil, err
	}

	messageID, err := ps.readID(r)
	if err != nil {
		return nil, err
	}

	content, err := ps.readByteArray(r)
	if err != nil {
		return nil, err
	}

	// Read headers map
	var headerCount uint32
	if err := binary.Read(r, binary.BigEndian, &headerCount); err != nil {
		return nil, err
	}

	headers := make(map[string]string, headerCount)
	for i := uint32(0); i < headerCount; i++ {
		key, err := ps.readString(r)
		if err != nil {
			return nil, err
		}
		value, err := ps.readString(r)
		if err != nil {
			return nil, err
		}
		headers[key] = value
	}

	return frames.CreateMessageFramePayload(header, topic, messageID, content, headers), nil
}

func (ps *Serializer) deserializeConnectPayload(r *bytes.Reader, header domain.HeaderPayload) (domain.ConnectFramePayload, error) {
	return nil, domain.ErrUnsupportedFrameType
}

func (ps *Serializer) readUint16(r *bytes.Reader) (uint16, error) {
	var val uint16
	err := binary.Read(r, binary.BigEndian, &val)
	return val, err
}

func (ps *Serializer) readUint32(r *bytes.Reader) (uint32, error) {
	var val uint32
	err := binary.Read(r, binary.BigEndian, &val)
	return val, err
}

func (ps *Serializer) readByteArray(r *bytes.Reader) ([]byte, error) {
	length, err := ps.readUint32(r)
	if err != nil {
		return nil, err
	}

	data := make([]byte, length)
	_, err = r.Read(data)
	return data, err
}

func (ps *Serializer) readString(r *bytes.Reader) (string, error) {
	data, err := ps.readByteArray(r)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (ps *Serializer) readID(r *bytes.Reader) (domain.ID, error) {
	str, err := ps.readString(r)
	if err != nil {
		return "", err
	}
	return domain.ID(str), nil
}

// NewSerializer return a new instance of serializer.
func NewSerializer(pool domain.Pool[*bytes.Buffer]) *Serializer {
	return &Serializer{
		bufferPool: pool,
	}
}
