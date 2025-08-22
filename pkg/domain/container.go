package domain

// ContainerState represent the current state of a container.
type ContainerState string

const (
	// CONTAINER_CREATED State represent the Container CONTAINER_CREATED Waiting for handshake process.
	CONTAINER_CREATED ContainerState = "CONTAINER_CREATED"
	// OPEN_SENT State represent the first phase of the handshake sending the OpenFrame to the client.
	OPEN_SENT ContainerState = "OPEN_SENT"
	// OPEN_RCVD State represent the confirmation of the received OpenFrame from the client.
	OPEN_RCVD ContainerState = "OPEN_RCVD"
	// CONNECTED State represent the validation and connection from the client to the broker and fully assigned to his container.
	CONNECTED ContainerState = "CLIENT_CON"
	// RESERVED State represent the state when a client shutdown and the container have no more clients reserve the container while clients are rebooting.
	RESERVED ContainerState = "RSRVD"

	// IDLE State represent the state when no clients have been CONNECTED for a while so could be stole by a new one (overriding topics and all).
	IDLE ContainerState = "IDLE"
)

type Container interface {
	CreateChannel(topic string, idGenerator func() ID) Channel
	RemoveChannel(topic string)
	SetState(state ContainerState)
}

type Channel interface {
}
