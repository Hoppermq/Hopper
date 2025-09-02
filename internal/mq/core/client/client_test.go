package client

import (
	"context"
	"testing"

	"github.com/hoppermq/hopper/internal/common"
	"github.com/hoppermq/hopper/pkg/domain"
	"github.com/hoppermq/hopper/pkg/domain/mocks"
	mocks_generator "github.com/hoppermq/hopper/pkg/domain/mocks/common"
	"github.com/stretchr/testify/assert"
)

func TestNewClientManager(t *testing.T) {
	t.Parallel()

	type args struct{}
	tests := []struct {
		name string
		args args
		want bool
	}{{
		name: "TestNewClientManager",
		args: args{},
		want: true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := NewManager(common.GenerateIdentifier)
			assert.Equal(t, tt.want, got != nil)
		})
	}
}

func TestClientManager_CreateNewClient(t *testing.T) {
	t.Parallel()

	type args struct {
		conn domain.Connection
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "CreateNewClient",
			args: args{
				conn: &mocks.MockConnection{},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cm := NewManager(common.GenerateIdentifier)
			got := cm.createClient(tt.args.conn)
			assert.Equal(t, tt.want, got != nil)
		})
	}
}

func TestClientManager_ClientHandling(t *testing.T) {
	t.Parallel()

	type args struct {
		conn       domain.Connection
		expectedID domain.ID
		setupMock  func(*mocks_generator.MockGenerator)
	}
	tests := []struct {
		name        string
		args        args
		wantIDCheck assert.ValueAssertionFunc
	}{{
		name: "ClientHandlingNewClient",
		args: args{
			conn:       mocks.NewMockConnection(t),
			expectedID: "test-client-id-12345",
			setupMock: func(g *mocks_generator.MockGenerator) {
				g.ON("test-client-id-12345")
			},
		},
		wantIDCheck: func(t assert.TestingT, client interface{}, i2 ...interface{}) bool {
			c, ok := client.(*Client)
			if !ok {
				return false
			}
			return assert.Equal(t, domain.ID("test-client-id-12345"), c.ID)
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockGen := mocks_generator.NewMockGenerator(common.GenerateIdentifier)
			tt.args.setupMock(mockGen)

			cm := NewManager(mockGen.Generator)
			client := cm.HandleNewClient(tt.args.conn)

			tt.wantIDCheck(t, client)

			assert.NotNil(t, client, "Client should not be nil")
			assert.Equal(t, tt.args.expectedID, client.ID, "Client should have the expected ID")
			assert.Equal(t, tt.args.conn, client.Conn, "Client should have the correct connection")

			retrievedClient := cm.GetClient(client.ID)
			assert.Equal(t, client, retrievedClient, "Client should be retrievable from manager")

			clientByConn := cm.GetClientByConnection(tt.args.conn)
			assert.Equal(t, client, clientByConn, "Client should be retrievable by connection")
		})
	}
}

func TestClientManager_RemoveClient(t *testing.T) {
	t.Parallel()
	type args struct {
		setupConnection func() domain.Connection
		setupMock       func(*mocks_generator.MockGenerator)
		clientID        domain.ID
	}

	tests := []struct {
		name        string
		args        args
		wantSuccess assert.ValueAssertionFunc
	}{
		{
			name: "RemoveExistingClient",
			args: args{
				setupConnection: func() domain.Connection {
					mockConn := mocks.NewMockConnection(t)
					mockConn.On("Close").Return(nil).Once()
					return mockConn
				},
				setupMock: func(g *mocks_generator.MockGenerator) {
					g.ON("test-client-remove-123")
				},
				clientID: "test-client-remove-123",
			},
			wantSuccess: func(t assert.TestingT, value interface{}, msgAndArgs ...interface{}) bool {
				args := value.(args)
				mockGen := mocks_generator.NewMockGenerator(common.GenerateIdentifier)
				args.setupMock(mockGen)
				cm := NewManager(mockGen.Generator)
				conn := args.setupConnection()

				client := cm.HandleNewClient(conn)
				assert.Equal(t, args.clientID, client.ID)
				assert.NotNil(t, cm.GetClient(client.ID), "Client should exist before removal")

				cm.RemoveClient(client.ID)

				return assert.Nil(t, cm.GetClient(client.ID), "Client should be nil after removal") &&
					assert.Nil(t, cm.GetClientByConnection(conn), "Client should not be findable by connection")
			},
		},
		{
			name: "RemoveNonExistentClient",
			args: args{
				setupConnection: func() domain.Connection {
					return mocks.NewMockConnection(t)
				},
				setupMock: func(g *mocks_generator.MockGenerator) {
					g.ON("existing-client-id")
				},
				clientID: "non-existent-id",
			},
			wantSuccess: func(t assert.TestingT, value interface{}, msgAndArgs ...interface{}) bool {
				args := value.(args)
				mockGen := mocks_generator.NewMockGenerator(common.GenerateIdentifier)
				args.setupMock(mockGen)
				cm := NewManager(mockGen.Generator)

				conn := args.setupConnection()
				client := cm.HandleNewClient(conn)
				assert.NotEqual(t, args.clientID, client.ID)

				cm.RemoveClient(domain.ID(args.clientID))

				return assert.NotNil(t, cm.GetClient(client.ID), "Existing client should still be present")
			},
		},
		{
			name: "RemoveClientWithConnectionCloseError",
			args: args{
				setupConnection: func() domain.Connection {
					mockConn := mocks.NewMockConnection(t)
					mockConn.On("Close").Return(assert.AnError).Once()
					return mockConn
				},
				setupMock: func(g *mocks_generator.MockGenerator) {
					g.ON("test-client-close-error")
				},
				clientID: "test-client-close-error",
			},
			wantSuccess: func(t assert.TestingT, value interface{}, msgAndArgs ...interface{}) bool {
				args := value.(args)
				mockGen := mocks_generator.NewMockGenerator(common.GenerateIdentifier)
				args.setupMock(mockGen)
				cm := NewManager(mockGen.Generator)
				conn := args.setupConnection()

				client := cm.HandleNewClient(conn)
				assert.Equal(t, args.clientID, client.ID)

				cm.RemoveClient(client.ID)

				return assert.NotNil(t, cm.GetClient(client.ID), "Client should still exist due to connection close error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.wantSuccess != nil {
				tt.wantSuccess(t, tt.args)
			}
		})
	}
}

func TestClientManager_GetClient(t *testing.T) {
	t.Parallel()
	type args struct {
		setupConnection func() domain.Connection
		setupMock       func(*mocks_generator.MockGenerator)
		targetClientID  domain.ID
	}

	tests := []struct {
		name        string
		args        args
		wantSuccess assert.ValueAssertionFunc
	}{
		{
			name: "GetExistingClient",
			args: args{
				setupConnection: func() domain.Connection {
					return mocks.NewMockConnection(t)
				},
				setupMock: func(g *mocks_generator.MockGenerator) {
					g.ON("get-existing-client-456")
				},
				targetClientID: "get-existing-client-456",
			},
			wantSuccess: func(t assert.TestingT, value interface{}, msgAndArgs ...interface{}) bool {
				args := value.(args)
				mockGen := mocks_generator.NewMockGenerator(common.GenerateIdentifier)
				args.setupMock(mockGen)
				cm := NewManager(mockGen.Generator)
				conn := args.setupConnection()

				client := cm.HandleNewClient(conn)
				assert.Equal(t, args.targetClientID, client.ID)

				retrievedClient := cm.GetClient(client.ID)

				return assert.NotNil(t, retrievedClient, "Retrieved client should not be nil") &&
					assert.Equal(t, client.ID, retrievedClient.ID, "Client IDs should match") &&
					assert.Equal(t, client.Conn, retrievedClient.Conn, "Client connections should match") &&
					assert.Equal(t, client, retrievedClient, "Clients should be identical")
			},
		},
		{
			name: "GetNonExistentClient",
			args: args{
				setupConnection: func() domain.Connection {
					return mocks.NewMockConnection(t)
				},
				setupMock: func(g *mocks_generator.MockGenerator) {
					g.ON("different-client-id")
				},
				targetClientID: "non-existent-client-id",
			},
			wantSuccess: func(t assert.TestingT, value interface{}, msgAndArgs ...interface{}) bool {
				args := value.(args)
				mockGen := mocks_generator.NewMockGenerator(common.GenerateIdentifier)
				args.setupMock(mockGen)
				cm := NewManager(mockGen.Generator)
				conn := args.setupConnection()

				client := cm.HandleNewClient(conn)
				assert.NotEqual(t, args.targetClientID, client.ID)

				retrievedClient := cm.GetClient(args.targetClientID)

				return assert.Nil(t, retrievedClient, "Non-existent client should return nil")
			},
		},
		{
			name: "GetClientFromEmptyManager",
			args: args{
				setupConnection: func() domain.Connection {
					return mocks.NewMockConnection(t)
				},
				setupMock: func(g *mocks_generator.MockGenerator) {
				},
				targetClientID: "any-client-id",
			},
			wantSuccess: func(t assert.TestingT, value interface{}, msgAndArgs ...interface{}) bool {
				args := value.(args)
				mockGen := mocks_generator.NewMockGenerator(common.GenerateIdentifier)
				args.setupMock(mockGen)
				cm := NewManager(mockGen.Generator)

				retrievedClient := cm.GetClient(args.targetClientID)

				return assert.Nil(t, retrievedClient, "Empty manager should return nil")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.wantSuccess != nil {
				tt.wantSuccess(t, tt.args)
			}
		})
	}
}

func TestClientManager_GetClientByConnection(t *testing.T) {
	t.Parallel()
	type args struct {
		setupConnections func() []domain.Connection
		setupMock        func(*mocks_generator.MockGenerator)
	}

	tests := []struct {
		name        string
		args        args
		wantSuccess assert.ValueAssertionFunc
	}{
		{
			name: "GetClientByExistingConnection",
			args: args{
				setupConnections: func() []domain.Connection {
					conn := mocks.NewMockConnection(t)
					return []domain.Connection{conn}
				},
				setupMock: func(g *mocks_generator.MockGenerator) {
					g.ON("conn-client-789")
				},
			},
			wantSuccess: func(t assert.TestingT, value interface{}, msgAndArgs ...interface{}) bool {
				args := value.(args)
				mockGen := mocks_generator.NewMockGenerator(common.GenerateIdentifier)
				args.setupMock(mockGen)
				cm := NewManager(mockGen.Generator)
				connections := args.setupConnections()
				conn := connections[0]

				client := cm.HandleNewClient(conn)

				retrievedClient := cm.GetClientByConnection(conn)

				return assert.NotNil(t, retrievedClient, "Client should be found by connection") &&
					assert.Equal(t, client, retrievedClient, "Retrieved client should match original") &&
					assert.Equal(t, conn, retrievedClient.Conn, "Connection should match")
			},
		},
		{
			name: "GetClientByNonExistentConnection",
			args: args{
				setupConnections: func() []domain.Connection {
					conn1 := mocks.NewMockConnection(t)
					conn2 := mocks.NewMockConnection(t) // This will be the non-existent one
					return []domain.Connection{conn1, conn2}
				},
				setupMock: func(g *mocks_generator.MockGenerator) {
					g.ON("conn-client-456")
				},
			},
			wantSuccess: func(t assert.TestingT, value interface{}, msgAndArgs ...interface{}) bool {
				args := value.(args)
				mockGen := mocks_generator.NewMockGenerator(common.GenerateIdentifier)
				args.setupMock(mockGen)
				cm := NewManager(mockGen.Generator)
				connections := args.setupConnections()
				conn1, conn2 := connections[0], connections[1]

				cm.HandleNewClient(conn1)

				retrievedClient := cm.GetClientByConnection(conn2)

				return assert.Nil(t, retrievedClient, "Should return nil for non-existent connection")
			},
		},
		{
			name: "GetClientByConnectionFromEmptyManager",
			args: args{
				setupConnections: func() []domain.Connection {
					return []domain.Connection{mocks.NewMockConnection(t)}
				},
				setupMock: func(g *mocks_generator.MockGenerator) {
					// No setup needed
				},
			},
			wantSuccess: func(t assert.TestingT, value interface{}, msgAndArgs ...interface{}) bool {
				args := value.(args)
				mockGen := mocks_generator.NewMockGenerator(common.GenerateIdentifier)
				args.setupMock(mockGen)
				cm := NewManager(mockGen.Generator)
				connections := args.setupConnections()

				// Try to get client from empty manager
				retrievedClient := cm.GetClientByConnection(connections[0])

				return assert.Nil(t, retrievedClient, "Empty manager should return nil")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.wantSuccess != nil {
				tt.wantSuccess(t, tt.args)
			}
		})
	}
}

func TestClientManager_Shutdown(t *testing.T) {
	t.Parallel()
	type args struct {
		setupConnections func() []domain.Connection
		setupMock        func(*mocks_generator.MockGenerator)
	}

	tests := []struct {
		name        string
		args        args
		wantSuccess assert.ValueAssertionFunc
	}{
		{
			name: "ShutdownWithMultipleClients",
			args: args{
				setupConnections: func() []domain.Connection {
					conn1 := mocks.NewMockConnection(t)
					conn2 := mocks.NewMockConnection(t)

					// Each connection should expect a Close() call during shutdown
					conn1.On("Close").Return(nil).Once()
					conn2.On("Close").Return(nil).Once()

					return []domain.Connection{conn1, conn2}
				},
				setupMock: func(g *mocks_generator.MockGenerator) {
					// Set up generator to cycle through different IDs
					call := 0
					g.Generator = func() domain.ID {
						call++
						return domain.ID("shutdown-client-" + string(rune('0'+call)))
					}
				},
			},
			wantSuccess: func(t assert.TestingT, value interface{}, msgAndArgs ...interface{}) bool {
				args := value.(args)
				mockGen := mocks_generator.NewMockGenerator(common.GenerateIdentifier)
				args.setupMock(mockGen)
				cm := NewManager(mockGen.Generator)
				connections := args.setupConnections()

				// Create multiple clients
				clients := make([]*Client, len(connections))
				for i, conn := range connections {
					clients[i] = cm.HandleNewClient(conn)
				}

				// Verify clients exist before shutdown
				for _, client := range clients {
					assert.NotNil(t, cm.GetClient(client.ID), "Client should exist before shutdown")
				}

				// Shutdown
				err := cm.Shutdown(context.Background())

				// Verify shutdown completed successfully
				return assert.NoError(t, err, "Shutdown should complete without error")
			},
		},
		{
			name: "ShutdownEmptyManager",
			args: args{
				setupConnections: func() []domain.Connection {
					return []domain.Connection{}
				},
				setupMock: func(g *mocks_generator.MockGenerator) {
					// No setup needed
				},
			},
			wantSuccess: func(t assert.TestingT, value interface{}, msgAndArgs ...interface{}) bool {
				args := value.(args)
				mockGen := mocks_generator.NewMockGenerator(common.GenerateIdentifier)
				args.setupMock(mockGen)
				cm := NewManager(mockGen.Generator)

				// Shutdown empty manager
				err := cm.Shutdown(context.Background())

				return assert.NoError(t, err, "Shutdown of empty manager should not error")
			},
		},
		{
			name: "ShutdownWithConnectionErrors",
			args: args{
				setupConnections: func() []domain.Connection {
					conn1 := mocks.NewMockConnection(t)
					conn2 := mocks.NewMockConnection(t)

					// First connection closes successfully, second fails
					conn1.On("Close").Return(nil).Once()
					conn2.On("Close").Return(assert.AnError).Once()

					return []domain.Connection{conn1, conn2}
				},
				setupMock: func(g *mocks_generator.MockGenerator) {
					call := 0
					g.Generator = func() domain.ID {
						call++
						return domain.ID("error-client-" + string(rune('0'+call)))
					}
				},
			},
			wantSuccess: func(t assert.TestingT, value interface{}, msgAndArgs ...interface{}) bool {
				args := value.(args)
				mockGen := mocks_generator.NewMockGenerator(common.GenerateIdentifier)
				args.setupMock(mockGen)
				cm := NewManager(mockGen.Generator)
				connections := args.setupConnections()

				// Create clients
				for _, conn := range connections {
					cm.HandleNewClient(conn)
				}

				// Shutdown (should handle connection close errors gracefully)
				err := cm.Shutdown(context.Background())

				return assert.NoError(t, err, "Shutdown should complete even with connection close errors")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.wantSuccess != nil {
				tt.wantSuccess(t, tt.args)
			}
		})
	}
}
