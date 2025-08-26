package core

import (
	"context"
	"sync"

	"github.com/hoppermq/hopper/pkg/domain"
)

// Client represents a single client connection to the broker.
type Client struct {
	ID   domain.ID
	Conn domain.Connection
	Mut  sync.Mutex
}

// ClientManager is responsible for managing client connections to the broker.
type ClientManager struct {
	client map[domain.ID]*Client
	mut    sync.RWMutex
}

// NewClientManager creates a new ClientManager instance with the provided Broker.
func NewClientManager(b *Broker) *ClientManager {
	return &ClientManager{
		client: make(map[domain.ID]*Client),
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
func (cm *ClientManager) RemoveClient(clientID domain.ID) {
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

// GetClientByConnection finds a client by their connection.
func (cm *ClientManager) GetClientByConnection(conn domain.Connection) *Client {
	cm.mut.RLock()
	defer cm.mut.RUnlock()

	for _, client := range cm.client {
		if client.Conn == conn {
			return client
		}
	}
	return nil
}

// GetClient finds a client by their ID.
func (cm *ClientManager) GetClient(clientID domain.ID) *Client {
	cm.mut.RLock()
	defer cm.mut.RUnlock()

	return cm.client[clientID]
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
