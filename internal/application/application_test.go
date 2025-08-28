package application

import (
	"context"
	"log/slog"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/hoppermq/hopper/internal/config"
	"github.com/hoppermq/hopper/pkg/domain"
	"github.com/hoppermq/hopper/pkg/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateNewApplication(t *testing.T) {
	t.Run("should_create_a_new_application", func(t *testing.T) {
		app := New()
		assert.NotNil(t, app)
		assert.NotNil(t, app.running)
		assert.NotNil(t, app.stop)
		assert.Empty(t, app.services)
	})

	t.Run("should_create_application_with_options", func(t *testing.T) {
		cfg := &config.Configuration{
			App: struct {
				Name        string `koanf:"name"`
				Version     string `koanf:"version"`
				ID          string `koanf:"id"`
				Description string `koanf:"description"`
			}{
				Name:        "hoppermq_test",
				Version:     "v0.0.1",
				ID:          "hppmq-v0.0.1@master-280825152010",
				Description: "simple description",
			},
		}
		logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
		mockEventBus := mocks.NewMockIEventBus(t)

		app := New(
			WithConfiguration(cfg),
			WithLogger(logger),
			WithEventBus(mockEventBus),
		)

		assert.NotNil(t, app)
		assert.Equal(t, cfg, app.configuration)
		assert.Equal(t, logger, app.logger)
		assert.Equal(t, mockEventBus, app.eb)
	})
}

func TestApplication_Start(t *testing.T) {
	cfg := &config.Configuration{
		App: struct {
			Name        string `koanf:"name"`
			Version     string `koanf:"version"`
			ID          string `koanf:"id"`
			Description string `koanf:"description"`
		}{
			Name:        "test_app",
			Version:     "v1.0.0",
			ID:          "test-id",
			Description: "test description",
		},
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	t.Run("should_start_and_stop_application_with_signal", func(t *testing.T) {
		mockEventBus := mocks.NewMockIEventBus(t)
		mockService := mocks.NewMockService(t)

		mockService.EXPECT().Run(mock.Anything).Return(nil).Once()
		mockService.EXPECT().Stop(mock.Anything).Return(nil).Once()
		mockService.EXPECT().Name().Return("test-service").Maybe() // For logging

		app := New(
			WithConfiguration(cfg),
			WithLogger(logger),
			WithEventBus(mockEventBus),
			WithService(mockService),
		)

		done := make(chan bool)
		go func() {
			app.Start()
			done <- true
		}()

		time.Sleep(100 * time.Millisecond)

		app.stop <- syscall.SIGTERM

		select {
		case <-done:
			// Success
		case <-time.After(5 * time.Second):
			t.Fatal("Application didn't stop within timeout")
		}

		mockService.AssertExpectations(t)
	})

	t.Run("should_register_event_bus_for_eventbus_aware_services", func(t *testing.T) {
		mockEventBus := mocks.NewMockIEventBus(t)

		mockService := mocks.NewMockService(t)
		mockEventBusAware := mocks.NewMockEventBusAware(t)

		mockEventBusAware.EXPECT().RegisterEventBus(mockEventBus).Once()
		mockService.EXPECT().Run(mock.Anything).Return(nil).Once()
		mockService.EXPECT().Stop(mock.Anything).Return(nil).Once()
		mockService.EXPECT().Name().Return("eventbus-aware-service").Maybe()

		compositeService := &testEventBusAwareService{
			Service:       mockService,
			EventBusAware: mockEventBusAware,
		}

		app := New(
			WithConfiguration(cfg),
			WithLogger(logger),
			WithEventBus(mockEventBus),
			WithService(compositeService),
		)

		done := make(chan bool)
		go func() {
			app.Start()
			done <- true
		}()

		time.Sleep(100 * time.Millisecond)
		app.stop <- syscall.SIGTERM

		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("Application didn't stop within timeout")
		}

		mockService.AssertExpectations(t)
		mockEventBusAware.AssertExpectations(t)
	})

	t.Run("should_handle_service_startup_failure", func(t *testing.T) {
		mockEventBus := mocks.NewMockIEventBus(t)
		mockService := mocks.NewMockService(t)

		mockService.EXPECT().Run(mock.Anything).Return(assert.AnError).Once()
		mockService.EXPECT().Stop(mock.Anything).Return(nil).Once()
		mockService.EXPECT().Name().Return("failing-service").Maybe()

		app := New(
			WithConfiguration(cfg),
			WithLogger(logger),
			WithEventBus(mockEventBus),
			WithService(mockService),
		)

		done := make(chan bool)
		go func() {
			app.Start()
			done <- true
		}()

		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("Application didn't stop within timeout")
		}

		mockService.AssertExpectations(t)
	})
}

func TestApplication_Stop(t *testing.T) {
	cfg := &config.Configuration{
		App: struct {
			Name        string `koanf:"name"`
			Version     string `koanf:"version"`
			ID          string `koanf:"id"`
			Description string `koanf:"description"`
		}{
			Name: "test_app",
			ID:   "test-id",
		},
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	t.Run("should_stop_services_with_timeout", func(t *testing.T) {
		mockEventBus := mocks.NewMockIEventBus(t)
		mockService := mocks.NewMockService(t)

		mockService.EXPECT().Stop(mock.Anything).Return(nil).Once()
		mockService.EXPECT().Name().Return("test-service").Maybe()

		app := New(
			WithConfiguration(cfg),
			WithLogger(logger),
			WithEventBus(mockEventBus),
			WithService(mockService),
		)

		done := make(chan bool)
		go func() {
			app.Stop()
			done <- true
		}()

		time.Sleep(50 * time.Millisecond)
		app.stop <- syscall.SIGTERM

		select {
		case <-done:
		case <-time.After(35 * time.Second):
			t.Fatal("Stop didn't complete within timeout")
		}

		mockService.AssertExpectations(t)
	})
}

func TestApplication_getName(t *testing.T) {
	cfg := &config.Configuration{
		App: struct {
			Name        string `koanf:"name"`
			Version     string `koanf:"version"`
			ID          string `koanf:"id"`
			Description string `koanf:"description"`
		}{
			Name: "test-name",
		},
	}

	app := New(WithConfiguration(cfg))
	assert.Equal(t, "test-name", app.getName())
}

func TestApplication_getID(t *testing.T) {
	cfg := &config.Configuration{
		App: struct {
			Name        string `koanf:"name"`
			Version     string `koanf:"version"`
			ID          string `koanf:"id"`
			Description string `koanf:"description"`
		}{
			ID: "test-id-123",
		},
	}

	app := New(WithConfiguration(cfg))
	assert.Equal(t, domain.ID("test-id-123"), app.getID())
}

type testEventBusAwareService struct {
	domain.Service
	domain.EventBusAware
}

func (t *testEventBusAwareService) Name() string {
	return t.Service.Name()
}

func (t *testEventBusAwareService) Run(ctx context.Context) error {
	return t.Service.Run(ctx)
}

func (t *testEventBusAwareService) Stop(ctx context.Context) error {
	return t.Service.Stop(ctx)
}

func (t *testEventBusAwareService) RegisterEventBus(eventBus domain.IEventBus) {
	t.EventBusAware.RegisterEventBus(eventBus)
}
