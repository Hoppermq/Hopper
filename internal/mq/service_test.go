package mq

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/hoppermq/hopper/pkg/domain"
	"github.com/hoppermq/hopper/pkg/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateService(t *testing.T) {
	t.Parallel()
	type args struct {
		opts []Option
	}

	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr assert.ValueAssertionFunc
	}{
		{
			name: "Create_Service_With_Logger_And_TCP",
			args: args{
				opts: []Option{
					WithLogger(
						slog.New(
							slog.NewJSONHandler(os.Stdout, nil),
						),
					),
				},
			},
			want: true,
			wantErr: func(t assert.TestingT, value interface{}, msgAndArgs ...interface{}) bool {
				service := value.(*HopperMQService)
				return assert.Equal(t, "hopper-mq", service.Name()) &&
					assert.NotNil(t, service.logger) &&
					assert.NotNil(t, service.broker)
			},
		},
		{
			name: "Create_Service_With_Logger_Only",
			args: args{
				opts: []Option{
					WithLogger(
						slog.New(
							slog.NewJSONHandler(os.Stdout, nil),
						),
					),
				},
			},
			want: true,
			wantErr: func(t assert.TestingT, value interface{}, msgAndArgs ...interface{}) bool {
				service := value.(*HopperMQService)
				return assert.Equal(t, "hopper-mq", service.Name()) &&
					assert.NotNil(t, service.logger) &&
					assert.NotNil(t, service.broker) &&
					assert.Nil(t, service.tcpHandler)
			},
		},
		{
			name: "Create_Service_With_No_Options",
			args: args{
				opts: []Option{},
			},
			want: true,
			wantErr: func(t assert.TestingT, value interface{}, msgAndArgs ...interface{}) bool {
				service := value.(*HopperMQService)
				return assert.Equal(t, "hopper-mq", service.Name()) &&
					assert.NotNil(t, service.broker) &&
					assert.Nil(t, service.logger) &&
					assert.Nil(t, service.tcpHandler)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := New(tt.args.opts...)
			assert.Equal(t, tt.want, got != nil)

			if tt.wantErr != nil {
				tt.wantErr(t, got)
			}
		})
	}
}

func TestHopperMQService_RegisterEventBus(t *testing.T) {
	t.Parallel()
	type args struct {
		service *HopperMQService
		eb      domain.IEventBus
	}

	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr assert.ValueAssertionFunc
	}{
		{
			name: "RegisterEventBus_With_All_Components",
			args: args{
				service: func() *HopperMQService {
					mockBroker := mocks.NewMockIService(t)
					mockTransport := mocks.NewMockTransport(t)
					return &HopperMQService{
						broker:     mockBroker,
						tcpHandler: mockTransport,
					}
				}(),
				eb: mocks.NewMockIEventBus(t),
			},
			want: true,
			wantErr: func(t assert.TestingT, value interface{}, msgAndArgs ...interface{}) bool {
				args := value.(args)
				mockBroker := args.service.broker.(*mocks.MockIService)
				mockTransport := args.service.tcpHandler.(*mocks.MockTransport)
				mockEB := args.eb.(*mocks.MockIEventBus)

				mockBroker.EXPECT().RegisterEventBus(mockEB).Once()
				mockTransport.EXPECT().RegisterEventBus(mockEB).Once()

				args.service.RegisterEventBus(args.eb)
				return assert.Equal(t, args.eb, args.service.eb)
			},
		},
		{
			name: "RegisterEventBus_With_Nil_Components",
			args: args{
				service: &HopperMQService{
					broker:     nil,
					tcpHandler: nil,
				},
				eb: mocks.NewMockIEventBus(t),
			},
			want: true,
			wantErr: func(t assert.TestingT, value interface{}, msgAndArgs ...interface{}) bool {
				args := value.(args)
				args.service.RegisterEventBus(args.eb)
				return assert.Equal(t, args.eb, args.service.eb)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, tt.args.service != nil)

			if tt.wantErr != nil {
				tt.wantErr(t, tt.args)
			}
		})
	}
}

func TestHopperMQService_RunAndStop(t *testing.T) {
	t.Parallel()
	type args struct {
		service *HopperMQService
		ctx     context.Context
	}

	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr assert.ValueAssertionFunc
	}{
		{
			name: "Run_And_Stop_Successfully",
			args: args{
				service: func() *HopperMQService {
					mockBroker := mocks.NewMockIService(t)
					mockTransport := mocks.NewMockTransport(t)
					mockEB := mocks.NewMockIEventBus(t)
					logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
					return &HopperMQService{
						broker:     mockBroker,
						tcpHandler: mockTransport,
						eb:         mockEB,
						logger:     logger,
					}
				}(),
				ctx: context.Background(),
			},
			want: true,
			wantErr: func(t assert.TestingT, value interface{}, msgAndArgs ...interface{}) bool {
				args := value.(args)
				mockBroker := args.service.broker.(*mocks.MockIService)
				mockTransport := args.service.tcpHandler.(*mocks.MockTransport)

				mockBroker.EXPECT().Run(mock.Anything).Return(nil).Once()
				mockTransport.EXPECT().Stop(mock.Anything).Return(nil).Once()
				mockBroker.EXPECT().Stop(mock.Anything).Return(nil).Once()

				done := make(chan error, 1)
				go func() {
					done <- args.service.Run(args.ctx)
				}()

				time.Sleep(10 * time.Millisecond)

				stopErr := args.service.Stop(context.Background())
				assert.NoError(t, stopErr)

				select {
				case runErr := <-done:
					return assert.NoError(t, runErr)
				case <-time.After(200 * time.Millisecond):
					t.Errorf("Run did not complete in expected time")
					return false
				}
			},
		},
		{
			name: "Run_With_Broker_Error",
			args: args{
				service: func() *HopperMQService {
					mockBroker := mocks.NewMockIService(t)
					mockTransport := mocks.NewMockTransport(t)
					mockEB := mocks.NewMockIEventBus(t)
					logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
					return &HopperMQService{
						broker:     mockBroker,
						tcpHandler: mockTransport,
						eb:         mockEB,
						logger:     logger,
					}
				}(),
				ctx: context.Background(),
			},
			want: true,
			wantErr: func(t assert.TestingT, value interface{}, msgAndArgs ...interface{}) bool {
				args := value.(args)
				mockBroker := args.service.broker.(*mocks.MockIService)
				mockTransport := args.service.tcpHandler.(*mocks.MockTransport)

				mockBroker.EXPECT().Run(mock.Anything).Return(errors.New("broker failed")).Once()
				mockTransport.EXPECT().Stop(mock.Anything).Return(nil).Once()
				mockBroker.EXPECT().Stop(mock.Anything).Return(nil).Once()

				done := make(chan error, 1)
				go func() {
					done <- args.service.Run(args.ctx)
				}()

				time.Sleep(20 * time.Millisecond)

				stopErr := args.service.Stop(context.Background())
				assert.NoError(t, stopErr)

				select {
				case runErr := <-done:
					return assert.NoError(t, runErr)
				case <-time.After(200 * time.Millisecond):
					t.Errorf("Run did not complete in expected time")
					return false
				}
			},
		},
		{
			name: "Stop_With_Component_Errors",
			args: args{
				service: func() *HopperMQService {
					mockBroker := mocks.NewMockIService(t)
					mockTransport := mocks.NewMockTransport(t)
					logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
					return &HopperMQService{
						broker:     mockBroker,
						tcpHandler: mockTransport,
						logger:     logger,
					}
				}(),
				ctx: context.Background(),
			},
			want: true,
			wantErr: func(t assert.TestingT, value interface{}, msgAndArgs ...interface{}) bool {
				args := value.(args)
				mockBroker := args.service.broker.(*mocks.MockIService)
				mockTransport := args.service.tcpHandler.(*mocks.MockTransport)

				mockTransport.EXPECT().Stop(mock.Anything).Return(errors.New("transport stop failed")).Once()
				mockBroker.EXPECT().Stop(mock.Anything).Return(errors.New("broker stop failed")).Once()

				stopErr := args.service.Stop(args.ctx)
				return assert.NoError(t, stopErr)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, tt.args.service != nil)

			if tt.wantErr != nil {
				tt.wantErr(t, tt.args)
			}
		})
	}
}
