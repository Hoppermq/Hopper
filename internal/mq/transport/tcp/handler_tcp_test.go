package tcp

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"os"
	"testing"

	"github.com/hoppermq/hopper/internal/events"
	"github.com/hoppermq/hopper/pkg/domain"
	"github.com/hoppermq/hopper/pkg/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewTCP(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx  context.Context
		opts []Option
	}

	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr assert.ValueAssertionFunc
	}{
		{
			name: "Create_TCP_Handler_With_Option_Error",
			args: args{
				ctx: context.Background(),
				opts: []Option{
					func(c *config) error {
						return errors.New("option error")
					},
				},
			},
			want: false,
			wantErr: func(t assert.TestingT, value interface{}, msgAndArgs ...interface{}) bool {
				err := value.(error)
				return assert.Error(t, err) &&
					assert.Contains(t, err.Error(), "option error")
			},
		},
		{
			name: "Create_TCP_Handler_With_Valid_Options",
			args: args{
				ctx: context.Background(),
				opts: []Option{
					WithLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil))),
					WithListener(&net.ListenConfig{}),
				},
			},
			want: true,
			wantErr: func(t assert.TestingT, value interface{}, msgAndArgs ...interface{}) bool {
				tcp := value.(*TCP)
				return assert.NotNil(t, tcp) &&
					assert.NotNil(t, tcp.logger) &&
					assert.Equal(t, "tcp", tcp.Name())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewTCP(tt.args.ctx, tt.args.opts...)

			if tt.want {
				assert.NoError(t, err)
				assert.NotNil(t, got)
			} else {
				if tt.wantErr != nil {
					tt.wantErr(t, err)
				}
				if err == nil {
					assert.NotNil(t, got)
					assert.Equal(t, "tcp-handler", got.Name())
				}
			}
		})
	}
}

func TestTCP_Options(t *testing.T) {
	t.Parallel()
	type args struct {
		option Option
	}

	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr assert.ValueAssertionFunc
	}{
		{
			name: "WithLogger_Option",
			args: args{
				option: WithLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil))),
			},
			want: true,
			wantErr: func(t assert.TestingT, value interface{}, msgAndArgs ...interface{}) bool {
				config := value.(*config)
				return assert.NotNil(t, config.logger)
			},
		},
		{
			name: "WithListener_Option",
			args: args{
				option: WithListener(&net.ListenConfig{}),
			},
			want: true,
			wantErr: func(t assert.TestingT, value interface{}, msgAndArgs ...interface{}) bool {
				config := value.(*config)
				return assert.NotNil(t, config.lconf)
			},
		},
		{
			name: "WithLogger_Nil_Option",
			args: args{
				option: WithLogger(nil),
			},
			want: true,
			wantErr: func(t assert.TestingT, value interface{}, msgAndArgs ...interface{}) bool {
				config := value.(*config)
				return assert.Nil(t, config.logger)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			config := &config{}
			err := tt.args.option(config)

			assert.Equal(t, tt.want, err == nil)

			if tt.wantErr != nil {
				tt.wantErr(t, config)
			}
		})
	}
}

func TestTCP_RegisterEventBus(t *testing.T) {
	t.Parallel()
	type args struct {
		tcp *TCP
		eb  domain.IEventBus
	}

	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr assert.ValueAssertionFunc
	}{
		{
			name: "RegisterEventBus_Successfully",
			args: args{
				tcp: &TCP{
					logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
				},
				eb: mocks.NewMockIEventBus(t),
			},
			want: true,
			wantErr: func(t assert.TestingT, value interface{}, msgAndArgs ...interface{}) bool {
				args := value.(args)
				args.tcp.RegisterEventBus(args.eb)
				return assert.Equal(t, args.eb, args.tcp.eb)
			},
		},
		{
			name: "RegisterEventBus_With_Valid_Logger",
			args: args{
				tcp: &TCP{
					logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
				},
				eb: mocks.NewMockIEventBus(t),
			},
			want: true,
			wantErr: func(t assert.TestingT, value interface{}, msgAndArgs ...interface{}) bool {
				args := value.(args)
				args.tcp.RegisterEventBus(args.eb)
				return assert.Equal(t, args.eb, args.tcp.eb)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, tt.args.tcp != nil)

			if tt.wantErr != nil {
				tt.wantErr(t, tt.args)
			}
		})
	}
}

func TestTCP_Name(t *testing.T) {
	t.Parallel()

	tcp := &TCP{}
	assert.Equal(t, "tcp-handler", tcp.Name())
}

func TestTCP_Stop(t *testing.T) {
	t.Parallel()
	type args struct {
		tcp *TCP
		ctx context.Context
	}

	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr assert.ValueAssertionFunc
	}{
		{
			name: "Stop_Without_Run",
			args: args{
				tcp: &TCP{
					Listener: &mockListener{
						acceptCh: make(chan acceptResult, 1),
						closeCh:  make(chan struct{}),
					},
					logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
					eb:     mocks.NewMockIEventBus(t),
				},
				ctx: context.Background(),
			},
			want: true,
			wantErr: func(t assert.TestingT, value interface{}, msgAndArgs ...interface{}) bool {
				args := value.(args)

				stopErr := args.tcp.Stop(args.ctx)
				return assert.NoError(t, stopErr)
			},
		},
		{
			name: "Stop_With_Nil_Cancel",
			args: args{
				tcp: &TCP{
					Listener: &mockListener{
						acceptCh: make(chan acceptResult, 1),
						closeCh:  make(chan struct{}),
					},
					logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
					cancel: nil,
				},
				ctx: context.Background(),
			},
			want: true,
			wantErr: func(t assert.TestingT, value interface{}, msgAndArgs ...interface{}) bool {
				args := value.(args)

				stopErr := args.tcp.Stop(args.ctx)
				return assert.NoError(t, stopErr)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, tt.args.tcp != nil)

			if tt.wantErr != nil {
				tt.wantErr(t, tt.args)
			}
		})
	}
}

func TestTCP_ComponentValidation(t *testing.T) {
	t.Parallel()
	type args struct {
		tcp *TCP
	}

	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr assert.ValueAssertionFunc
	}{
		{
			name: "TCP_Component_Properties",
			args: args{
				tcp: &TCP{
					logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
					eb:     mocks.NewMockIEventBus(t),
					Listener: &mockListener{
						acceptCh: make(chan acceptResult, 1),
						closeCh:  make(chan struct{}),
					},
				},
			},
			want: true,
			wantErr: func(t assert.TestingT, value interface{}, msgAndArgs ...interface{}) bool {
				args := value.(args)

				return assert.Equal(t, "tcp-handler", args.tcp.Name()) &&
					assert.NotNil(t, args.tcp.logger) &&
					assert.NotNil(t, args.tcp.eb) &&
					assert.NotNil(t, args.tcp.Listener)
			},
		},
		{
			name: "TCP_Component_With_Minimal_Setup",
			args: args{
				tcp: &TCP{
					logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
				},
			},
			want: true,
			wantErr: func(t assert.TestingT, value interface{}, msgAndArgs ...interface{}) bool {
				args := value.(args)

				return assert.Equal(t, "tcp-handler", args.tcp.Name()) &&
					assert.NotNil(t, args.tcp.logger)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, tt.args.tcp != nil)

			if tt.wantErr != nil {
				tt.wantErr(t, tt.args)
			}
		})
	}
}

func TestTCP_SendMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		message   []byte
		setupMock func(*mocks.MockConnection)
	}{
		{
			name:    "SendMessage_Successfully",
			message: []byte("test message"),
			setupMock: func(mockConn *mocks.MockConnection) {
				mockConn.On("SetWriteDeadline", mock.AnythingOfType("time.Time")).Return(nil).Once()
				mockConn.On("Write", []byte("test message")).Return(12, nil).Once()
			},
		},
		{
			name:    "SendMessage_With_Write_Error",
			message: []byte("failed message"),
			setupMock: func(mockConn *mocks.MockConnection) {
				mockConn.On("SetWriteDeadline", mock.AnythingOfType("time.Time")).Return(nil).Once()
				mockConn.On("Write", []byte("failed message")).Return(0, errors.New("write failed")).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tcp := &TCP{
				logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
			}

			mockConn := mocks.NewMockConnection(t)
			tt.setupMock(mockConn)

			tcp.sendMessage(tt.message, mockConn)
		})
	}
}

func TestTCP_HandleMessageSending(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		message   []byte
		setupMock func(*mocks.MockConnection)
		canceled  bool
	}{
		{
			name:    "HandleMessageSending_With_SendMessageEvent",
			message: []byte("test message"),
			setupMock: func(mockConn *mocks.MockConnection) {
				mockConn.On("SetWriteDeadline", mock.AnythingOfType("time.Time")).Return(nil).Once()
				mockConn.On("Write", []byte("test message")).Return(12, nil).Once()
			},
			canceled: false,
		},
		{
			name:      "HandleMessageSending_With_Context_Cancel",
			message:   []byte("canceled message"),
			setupMock: func(mockConn *mocks.MockConnection) {},
			canceled:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tcp := &TCP{
				logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
			}

			ctx := context.Background()
			if tt.canceled {
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(context.Background())
				cancel()
			}

			ch := make(chan domain.Event, 1)
			mockConn := mocks.NewMockConnection(t)

			tt.setupMock(mockConn)

			if !tt.canceled {
				ch <- &events.SendMessageEvent{
					ClientID:  "test-client",
					Conn:      mockConn,
					Message:   tt.message,
					Transport: domain.TransportTypeTCP,
					BaseEvent: events.BaseEvent{
						EventType: domain.EventTypeSendMessage,
					},
				}
			}
			close(ch)

			tcp.handleMessageSending(ctx, ch)
		})
	}
}

func TestTCP_ProcessConnection(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(*mocks.MockConnection, *mocks.MockIEventBus)
	}{
		{
			name: "ProcessConnection_Handles_Connection_Lifecycle_Correctly",
			setupMocks: func(mockConn *mocks.MockConnection, mockEB *mocks.MockIEventBus) {
				mockEB.On("Publish", mock.Anything, mock.MatchedBy(func(evt domain.Event) bool {
					newConnEvt, ok := evt.(*events.NewConnectionEvent)
					return ok && newConnEvt.Transport == domain.TransportTypeTCP
				})).Return(nil).Once()

				mockConn.On("SetReadDeadline", mock.AnythingOfType("time.Time")).Return(nil).Once()
				mockConn.On("Read", mock.AnythingOfType("[]uint8")).Return(0, io.EOF).Once()

				mockEB.On("Publish", mock.Anything, mock.MatchedBy(func(evt domain.Event) bool {
					disconnEvt, ok := evt.(*events.ClientDisconnectedEvent)
					return ok && disconnEvt.Transport == domain.TransportTypeTCP
				})).Return(nil).Once()

				mockConn.On("Close").Return(nil).Once()
			},
		},
		{
			name: "ProcessConnection_Handles_Successful_Message_Reading",
			setupMocks: func(mockConn *mocks.MockConnection, mockEB *mocks.MockIEventBus) {
				mockEB.On("Publish", mock.Anything, mock.MatchedBy(func(evt domain.Event) bool {
					newConnEvt, ok := evt.(*events.NewConnectionEvent)
					return ok && newConnEvt.Transport == domain.TransportTypeTCP
				})).Return(nil).Once()

				mockConn.On("SetReadDeadline", mock.AnythingOfType("time.Time")).Return(nil).Once()
				mockConn.On("Read", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
					buf := args[0].([]byte)
					data := []byte("hello world\n")
					copy(buf, data)
				}).Return(len("hello world\n"), nil).Once()

				mockEB.On("Publish", mock.Anything, mock.MatchedBy(func(evt domain.Event) bool {
					msgEvt, ok := evt.(*events.MessageReceivedEvent)
					return ok &&
						msgEvt.Transport == domain.TransportTypeTCP &&
						string(msgEvt.Message) == "hello world\n"
				})).Return(nil).Once()

				mockConn.On("SetReadDeadline", mock.AnythingOfType("time.Time")).Return(nil).Once()
				mockConn.On("Read", mock.AnythingOfType("[]uint8")).Return(0, io.EOF).Once()

				mockEB.On("Publish", mock.Anything, mock.MatchedBy(func(evt domain.Event) bool {
					disconnEvt, ok := evt.(*events.ClientDisconnectedEvent)
					return ok && disconnEvt.Transport == domain.TransportTypeTCP
				})).Return(nil).Once()

				mockConn.On("Close").Return(nil).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockConn := mocks.NewMockConnection(t)
			mockEB := mocks.NewMockIEventBus(t)

			tcp := &TCP{
				logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
				eb:     mockEB,
			}

			tt.setupMocks(mockConn, mockEB)
			tcp.processConnection(mockConn, context.Background())
		})
	}
}

func TestTCP_ReceiveMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setupMocks  func(*mocks.MockConnection, *mocks.MockIEventBus)
		expectError bool
	}{
		{
			name: "ReceiveMessage_Successfully",
			setupMocks: func(mockConn *mocks.MockConnection, mockEB *mocks.MockIEventBus) {
				mockConn.On("SetReadDeadline", mock.AnythingOfType("time.Time")).Return(nil).Once()
				mockConn.On("Read", mock.AnythingOfType("[]uint8")).Run(func(args mock.Arguments) {
					buf := args[0].([]byte)
					data := []byte("hello world\n")
					copy(buf, data)
				}).Return(len("hello world\n"), nil).Once()

				mockEB.On("Publish", mock.Anything, mock.MatchedBy(func(evt domain.Event) bool {
					msgEvt, ok := evt.(*events.MessageReceivedEvent)
					return ok &&
						msgEvt.Transport == domain.TransportTypeTCP &&
						string(msgEvt.Message) == "hello world\n"
				})).Return(nil).Once()
			},
			expectError: false,
		},
		{
			name: "ReceiveMessage_With_EOF",
			setupMocks: func(mockConn *mocks.MockConnection, mockEB *mocks.MockIEventBus) {
				mockConn.On("SetReadDeadline", mock.AnythingOfType("time.Time")).Return(nil).Once()
				mockConn.On("Read", mock.AnythingOfType("[]uint8")).Return(0, io.EOF).Once()

				mockEB.On("Publish", mock.Anything, mock.MatchedBy(func(evt domain.Event) bool {
					disconnEvt, ok := evt.(*events.ClientDisconnectedEvent)
					return ok && disconnEvt.Transport == domain.TransportTypeTCP
				})).Return(nil).Once()
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockConn := mocks.NewMockConnection(t)
			mockEB := mocks.NewMockIEventBus(t)

			tcp := &TCP{
				logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
				eb:     mockEB,
			}

			tt.setupMocks(mockConn, mockEB)
			err := tcp.receiveMsg(mockConn, context.Background())
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

type mockListener struct {
	acceptCh chan acceptResult
	closeCh  chan struct{}
	closed   bool
}

type acceptResult struct {
	conn net.Conn
	err  error
}

func (m *mockListener) Accept() (net.Conn, error) {
	select {
	case result := <-m.acceptCh:
		return result.conn, result.err
	case <-m.closeCh:
		return nil, errors.New("listener closed")
	}
}

func (m *mockListener) Close() error {
	if !m.closed {
		close(m.closeCh)
		m.closed = true
	}
	return nil
}

func (m *mockListener) Addr() net.Addr {
	return &net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 9091,
	}
}
