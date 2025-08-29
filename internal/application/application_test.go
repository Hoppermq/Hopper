package application

import (
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/hoppermq/hopper/internal/config"
	"github.com/hoppermq/hopper/internal/events"
	"github.com/hoppermq/hopper/pkg/domain"
	"github.com/hoppermq/hopper/pkg/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateApplication(t *testing.T) {
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
			name: "Create_Application_With_Logger",
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
				app := value.(*Application)
				return assert.NotNil(t, app.logger) &&
					assert.NotNil(t, app.running) &&
					assert.NotNil(t, app.stop) &&
					assert.Len(t, app.services, 0)
			},
		},
		{
			name: "Create_Application_With_Configuration",
			args: args{
				opts: []Option{
					WithConfiguration(&config.Configuration{
						App: struct {
							Name        string `koanf:"name"`
							Version     string `koanf:"version"`
							ID          string `koanf:"id"`
							Description string `koanf:"description"`
						}{
							Name:    "test-app",
							Version: "1.0.0",
							ID:      "test-id",
						},
					}),
				},
			},
			want: true,
			wantErr: func(t assert.TestingT, value interface{}, msgAndArgs ...interface{}) bool {
				app := value.(*Application)
				return assert.NotNil(t, app.configuration) &&
					assert.Equal(t, "test-app", app.getName()) &&
					assert.Equal(t, domain.ID("test-id"), app.getID())
			},
		},
		{
			name: "Create_Application_With_All_Options",
			args: args{
				opts: []Option{
					WithLogger(
						slog.New(
							slog.NewJSONHandler(os.Stdout, nil),
						),
					),
					WithConfiguration(&config.Configuration{
						App: struct {
							Name        string `koanf:"name"`
							Version     string `koanf:"version"`
							ID          string `koanf:"id"`
							Description string `koanf:"description"`
						}{
							Name:    "hopper-test",
							Version: "0.1.0",
							ID:      "hopper-test-id",
						},
					}),
					WithEventBus(events.WithConfig(&config.Configuration{})),
				},
			},
			want: true,
			wantErr: func(t assert.TestingT, value interface{}, msgAndArgs ...interface{}) bool {
				app := value.(*Application)
				return assert.NotNil(t, app.logger) &&
					assert.NotNil(t, app.configuration) &&
					assert.NotNil(t, app.eb) &&
					assert.Equal(t, "hopper-test", app.getName()) &&
					assert.Equal(t, domain.ID("hopper-test-id"), app.getID())
			},
		},
		{
			name: "Create_Application_With_No_Options",
			args: args{
				opts: []Option{},
			},
			want: true,
			wantErr: func(t assert.TestingT, value interface{}, msgAndArgs ...interface{}) bool {
				app := value.(*Application)
				return assert.NotNil(t, app.running) &&
					assert.NotNil(t, app.stop) &&
					assert.Len(t, app.services, 0) &&
					assert.Nil(t, app.logger) &&
					assert.Nil(t, app.configuration) &&
					assert.Nil(t, app.eb)
			},
		},
		{
			name: "Create_Application_With_Services",
			args: args{
				opts: []Option{
					WithService(func() domain.IService {
						mockService := mocks.NewMockIService(t)
						mockService.EXPECT().Name().Return("test-service").Maybe()
						return mockService
					}()),
					WithLogger(
						slog.New(
							slog.NewJSONHandler(os.Stdout, nil),
						),
					),
				},
			},
			want: true,
			wantErr: func(t assert.TestingT, value interface{}, msgAndArgs ...interface{}) bool {
				app := value.(*Application)
				return assert.NotNil(t, app.logger) &&
					assert.Len(t, app.services, 1)
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

func TestApplication_ServiceLifecycle(t *testing.T) {
	t.Parallel()
	type args struct {
		app *Application
	}

	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr assert.ValueAssertionFunc
	}{
		{
			name: "Application_With_Single_Service_Lifecycle",
			args: args{
				app: func() *Application {
					mockService := mocks.NewMockIService(t)
					mockEB := mocks.NewMockIEventBus(t)
					cfg := &config.Configuration{
						App: struct {
							Name        string `koanf:"name"`
							Version     string `koanf:"version"`
							ID          string `koanf:"id"`
							Description string `koanf:"description"`
						}{
							Name: "test-app",
							ID:   "test-id",
						},
					}
					return &Application{
						configuration: cfg,
						logger:        slog.New(slog.NewJSONHandler(os.Stdout, nil)),
						services:      []domain.IService{mockService},
						eb:            mockEB,
						running:       make(chan bool, 1),
						stop:          make(chan os.Signal, 1),
					}
				}(),
			},
			want: true,
			wantErr: func(t assert.TestingT, value interface{}, msgAndArgs ...interface{}) bool {
				args := value.(args)
				app := args.app
				mockService := app.services[0].(*mocks.MockIService)
				mockEB := app.eb.(*mocks.MockIEventBus)

				mockService.EXPECT().RegisterEventBus(mockEB).Once()
				mockService.EXPECT().Run(mock.Anything).Return(nil).Once()
				mockService.EXPECT().Name().Return("test-service").Maybe()
				mockService.EXPECT().Stop(mock.Anything).Return(nil).Once()

				done := make(chan bool)
				go func() {
					defer func() { done <- true }()
					app.Start()
				}()

				time.Sleep(10 * time.Millisecond)

				app.stop <- os.Interrupt

				select {
				case <-done:
					return true
				case <-time.After(500 * time.Millisecond):
					t.Errorf("Application did not complete in expected time")
					return false
				}
			},
		},
		{
			name: "Application_With_Service_Error",
			args: args{
				app: func() *Application {
					mockService := mocks.NewMockIService(t)
					mockEB := mocks.NewMockIEventBus(t)
					cfg := &config.Configuration{
						App: struct {
							Name        string `koanf:"name"`
							Version     string `koanf:"version"`
							ID          string `koanf:"id"`
							Description string `koanf:"description"`
						}{
							Name: "error-app",
							ID:   "error-id",
						},
					}
					return &Application{
						configuration: cfg,
						logger:        slog.New(slog.NewJSONHandler(os.Stdout, nil)),
						services:      []domain.IService{mockService},
						eb:            mockEB,
						running:       make(chan bool, 1),
						stop:          make(chan os.Signal, 1),
					}
				}(),
			},
			want: true,
			wantErr: func(t assert.TestingT, value interface{}, msgAndArgs ...interface{}) bool {
				args := value.(args)
				app := args.app
				mockService := app.services[0].(*mocks.MockIService)
				mockEB := app.eb.(*mocks.MockIEventBus)

				mockService.EXPECT().RegisterEventBus(mockEB).Once()
				mockService.EXPECT().Run(mock.Anything).Return(errors.New("service failed")).Once()
				mockService.EXPECT().Name().Return("failing-service").Maybe()
				mockService.EXPECT().Stop(mock.Anything).Return(nil).Once()

				done := make(chan bool)
				go func() {
					defer func() { done <- true }()
					app.Start()
				}()

				select {
				case <-done:
					return true
				case <-time.After(500 * time.Millisecond):
					app.stop <- os.Interrupt
					<-done
					return true
				}
			},
		},
		{
			name: "Application_With_Multiple_Services",
			args: args{
				app: func() *Application {
					mockService1 := mocks.NewMockIService(t)
					mockService2 := mocks.NewMockIService(t)
					mockEB := mocks.NewMockIEventBus(t)
					cfg := &config.Configuration{
						App: struct {
							Name        string `koanf:"name"`
							Version     string `koanf:"version"`
							ID          string `koanf:"id"`
							Description string `koanf:"description"`
						}{
							Name: "multi-service-app",
							ID:   "multi-id",
						},
					}
					return &Application{
						configuration: cfg,
						logger:        slog.New(slog.NewJSONHandler(os.Stdout, nil)),
						services:      []domain.IService{mockService1, mockService2},
						eb:            mockEB,
						running:       make(chan bool, 1),
						stop:          make(chan os.Signal, 1),
					}
				}(),
			},
			want: true,
			wantErr: func(t assert.TestingT, value interface{}, msgAndArgs ...interface{}) bool {
				args := value.(args)
				app := args.app
				mockService1 := app.services[0].(*mocks.MockIService)
				mockService2 := app.services[1].(*mocks.MockIService)
				mockEB := app.eb.(*mocks.MockIEventBus)

				mockService1.EXPECT().RegisterEventBus(mockEB).Once()
				mockService1.EXPECT().Run(mock.Anything).Return(nil).Once()
				mockService1.EXPECT().Name().Return("service-1").Maybe()
				mockService1.EXPECT().Stop(mock.Anything).Return(nil).Once()

				mockService2.EXPECT().RegisterEventBus(mockEB).Once()
				mockService2.EXPECT().Run(mock.Anything).Return(nil).Once()
				mockService2.EXPECT().Name().Return("service-2").Maybe()
				mockService2.EXPECT().Stop(mock.Anything).Return(nil).Once()

				if !assert.Len(t, app.services, 2) {
					return false
				}

				done := make(chan bool)
				go func() {
					defer func() { done <- true }()
					app.Start()
				}()

				time.Sleep(10 * time.Millisecond)
				app.stop <- os.Interrupt

				select {
				case <-done:
					return true
				case <-time.After(500 * time.Millisecond):
					t.Errorf("Application did not complete in expected time")
					return false
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, tt.args.app != nil)

			if tt.wantErr != nil {
				tt.wantErr(t, tt.args)
			}
		})
	}
}
