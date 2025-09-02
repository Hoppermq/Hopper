package domain

// ContainerManager manages container lifecycle and creation based on business rules.
type ContainerManager interface {
	CreateNewContainer(idGenerator func() ID, clientID ID) Container
	RegisterContainerToTopic(topic string, containerID ID)
	RemoveContainerFromTopic(topic string, containerID ID)
}

// ClientManager manages client connections and lifecycle.
type ClientManager interface {
	HandleNewClient(conn Connection) Client
	GetClient(clientID ID) Client
	RemoveClient(clientID ID)
}

// FrameManager manages frames to categorize theme.
type FrameManager interface {
	IsControlFrame(f FrameType) bool
	IsMessageFrame(f FrameType) bool
	IsErrorFrame(f FrameType) bool
}

// Client represents a client connection (you might want to make this an interface too later).
type Client interface {
	GetID() ID
	GetConnection() Connection
	AttachContainer(containerID ID)
}
