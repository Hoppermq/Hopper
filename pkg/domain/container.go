package domain

// ContainerState represent the current state of a container.
type ContainerState string

const (
	// ContainerCreated State represent the Container ContainerCreated Waiting for handshake process.
	ContainerCreated ContainerState = "CONTAINER_CREATED"

	// OpenSent State represent the first phase of the handshake sending the OpenFrame to the client.
	OpenSent ContainerState = "OPEN_SENT"

	// OpenRcvd State represent the confirmation of the received OpenFrame from the client.
	OpenRcvd ContainerState = "OPEN_RCVD"

	// Connected State represent the validation and connection from the client
	// to the broker and fully assigned to his container.
	Connected ContainerState = "CLIENT_CON"

	// Reserved State represent the state when a client shutdown
	// and the container have no more clients reserve the container while clients are rebooting.
	Reserved ContainerState = "RSRVD"

	// Idle State represent the state when no clients have been CONNECTED
	// for a while so could be stolen by a new one (overriding topics and all).
	Idle ContainerState = "IDLE"
)

// Container represent an hopper container.
type Container interface {
	CreateChannel(topic string, generateIdentifier func() ID) *Channel
	RemoveChannel(topic string)
	SetState(state ContainerState)
	GetID() ID
}

// Channel represent the channel used by container.
type Channel interface {
	GetID() ID
}
