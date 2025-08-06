package core

import (
	"context"
	"net"
	"sync"
)

type Client struct {
	ID   string
	Conn net.Conn
	Mut  sync.Mutex
}

type ClientManager struct {
	client (map[string]*Client)
	broker *Broker
	mut    sync.RWMutex
}

func NewClientManager(b *Broker) *ClientManager {
	return &ClientManager{
		client: make(map[string]*Client),
		broker: b,
	}
}

func createClient(conn net.Conn) *Client {
	return &Client{
		ID:   "id",
		Conn: conn,
	}
}

func (cm *ClientManager) HandleNewClient(conn net.Conn) *Client {
	return createClient(conn)
}

func (cm *ClientManager) RemoveClient(clientID string) {
	cm.mut.Lock()
	defer cm.mut.Unlock()

	if client, exists := cm.client[clientID]; exists {
		client.Conn.Close()
		delete(cm.client, clientID)
		return
	}

	cm.broker.Logger.Warn("Client not found", "id", clientID)
}

func (cm *ClientManager) Shutdown(ctx context.Context) error {
	cm.mut.Lock()
	defer cm.mut.Unlock()

	cm.broker.Logger.Info("Disconnecting all clients", "count", len(cm.client))

	for id := range cm.client {
		cm.RemoveClient(id)
	}

	return nil
}
