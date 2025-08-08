package domain

// Serializable is an interface for serializable objects in the HopperMQ protocol.
type Serializable interface {
	Serialize() ([]byte, error)
	Deserialize(data []byte) (Serializable, error)
}

type Connection interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Close() error
}
