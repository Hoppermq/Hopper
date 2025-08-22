package core

import (
	"context"
	"sync"

	"github.com/hoppermq/hopper/pkg/domain"
)

// Client represents a single client connection to the broker.
type Client struct {
	ID   string
	Conn domain.Connection
	Mut  sync.Mutex
}

// ClientManager is responsible for managing client connections to the broker.
type ClientManager struct {
	client map[string]*Client
	mut    sync.RWMutex
}

// NewClientManager creates a new ClientManager instance with the provided Broker.
func NewClientManager(b *Broker) *ClientManager {
	return &ClientManager{
		client: make(map[string]*Client),
	}
}

func createClient(conn domain.Connection) *Client {
	return &Client{
		ID:   GenerateIdentifier(),
		Conn: conn,
	}
}

// HandleNewClient creates a new client connection and adds it to the ClientManager.
func (cm *ClientManager) HandleNewClient(conn domain.Connection) *Client {
	cm.mut.Lock()
	defer cm.mut.Unlock()

	client := createClient(conn)
	cm.client[client.ID] = client

	return client
}

// RemoveClient removes a client from the ClientManager by its ID.
func (cm *ClientManager) RemoveClient(clientID string) {
	cm.mut.Lock()
	defer cm.mut.Unlock()

	if client, exists := cm.client[clientID]; exists {
		err := client.Conn.Close()
		if err != nil {
			return
		}
		delete(cm.client, clientID)
		return
	}
}

// Shutdown gracefully disconnects all clients managed by the ClientManager.
func (cm *ClientManager) Shutdown(ctx context.Context) error {
	cm.mut.Lock()
	defer cm.mut.Unlock()

	for id := range cm.client {
		cm.RemoveClient(id)
	}

	return nil
}
