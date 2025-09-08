package client

import (
	"context"
	"sync"

	"github.com/hoppermq/hopper/pkg/domain"
)

// Manager is responsible for managing client connections to the broker.
type Manager struct {
	client    map[domain.ID]*Client
	generator domain.Generator
	mut       sync.RWMutex
}

// NewManager creates a new ClientManager instance with the provided generator.
func NewManager(generator domain.Generator) *Manager {
	return &Manager{
		client:    make(map[domain.ID]*Client),
		generator: generator,
	}
}

func (cm *Manager) createClient(conn domain.Connection) *Client {
	return &Client{
		ID:   cm.generator(),
		Conn: conn,
	}
}

// HandleNewClient creates a new client connection and adds it to the ClientManager.
func (cm *Manager) HandleNewClient(conn domain.Connection) *Client {
	cm.mut.Lock()
	defer cm.mut.Unlock()

	client := cm.createClient(conn)
	cm.client[client.ID] = client

	return client
}

// RemoveClient removes a client from the ClientManager by its ID.
func (cm *Manager) RemoveClient(clientID domain.ID) {
	cm.mut.Lock()
	defer cm.mut.Unlock()

	if client, exists := cm.client[clientID]; exists {
		client.Mut.Lock()
		if !client.closed && client.Conn != nil {
			err := client.Conn.Close()
			if err != nil {
				client.Mut.Unlock()
				return
			}
			client.closed = true
		}
		client.Mut.Unlock()

		delete(cm.client, clientID)
		return
	}
}

// GetClientByConnection finds a client by their connection.
func (cm *Manager) GetClientByConnection(conn domain.Connection) *Client {
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
func (cm *Manager) GetClient(clientID domain.ID) *Client {
	cm.mut.RLock()
	defer cm.mut.RUnlock()

	return cm.client[clientID]
}

// Shutdown gracefully disconnects all clients managed by the ClientManager.
func (cm *Manager) Shutdown(ctx context.Context) error {
	cm.mut.Lock()
	defer cm.mut.Unlock()

	for id, client := range cm.client {
		client.Mut.Lock()
		if !client.closed && client.Conn != nil {
			if err := client.Conn.Close(); err != nil {
				// Log the error but continue shutting down other connections
				// This ensures graceful shutdown even if some connections fail to close
			}

			client.closed = true
		}
		client.Mut.Unlock()

		delete(cm.client, id)
	}

	return nil
}
