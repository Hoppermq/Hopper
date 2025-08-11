package domain

import "time"

// Serializable is an interface for serializable objects in the HopperMQ protocol.
type Serializable interface {
	Serialize() ([]byte, error)
	Deserialize(data []byte) (Serializable, error)
}

type Serializer interface {
	SerializeFrame() ([]byte, error)
	DeserializeFrame(data []byte) (Frame, error)
}

// Connection is an interface that use the same functions as net/Conn.
type Connection interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Close() error

	SetDeadline(t time.Time) error
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
}
