package client

import (
	"sync"

	"github.com/hoppermq/hopper/pkg/domain"
)

// Client represents a single client connection to the broker.
type Client struct {
	ID          domain.ID
	containerID domain.ID
	Conn        domain.Connection
	Mut         sync.Mutex

	closed bool
}

// GetID return the client ID.
func (c *Client) GetID() domain.ID {
	return c.ID
}

// GetConn return the client Connection.
func (c *Client) GetConn() domain.Connection {
	return c.Conn
}

// GetContainer return the container attached to the client.
func (c *Client) GetContainer() domain.ID {
	return c.containerID
}
