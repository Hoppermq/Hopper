package container

import (
	"context"
	"testing"

	"github.com/hoppermq/hopper/pkg/domain"
	"github.com/hoppermq/hopper/pkg/domain/mocks"
)

func TestContainer_HandleConnectFrame(t *testing.T) {
	t.Run("HandleConnectFrame_ValidState_Success", func(t *testing.T) {
		container := NewContainer("container123", "client123")
		container.State = domain.ContainerOpenSent

		payload := mocks.NewMockConnectFramePayload(t)
		payload.On("GetSourceID").Return(domain.ID("client123")).Twice()

		mockFrame := mocks.NewMockFrame(t)
		mockFrame.On("GetPayload").Return(payload)

		callbackCount := 0
		testCallback := func(ctx context.Context, frame domain.Frame, clientID domain.ID) error {
			callbackCount++
			if frame.GetType() != domain.FrameTypeBegin {
				t.Errorf("Expected Begin frame, got %v", frame.GetType())
			}
			if clientID != "client123" {
				t.Errorf("Expected clientID 'client123', got %v", clientID)
			}
			return nil
		}

		ctx := context.Background()
		err := container.HandleConnectFrame(ctx, mockFrame, testCallback)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if container.State != domain.ContainerConnected {
			t.Errorf("Expected state %v, got %v", domain.ContainerConnected, container.State)
		}
		if callbackCount != 1 {
			t.Errorf("Expected 1 callback call, got %d", callbackCount)
		}
		if _, exists := container.ChannelsByTopic["__temp__"]; !exists {
			t.Error("Expected temporary channel to be created")
		}
	})

	t.Run("HandleConnectFrame_InvalidState_Error", func(t *testing.T) {
		container := NewContainer("container123", "client123")
		container.State = domain.ContainerCreated

		mockFrame := mocks.NewMockFrame(t)

		callbackCount := 0
		testCallback := func(ctx context.Context, frame domain.Frame, clientID domain.ID) error {
			callbackCount++
			return nil
		}

		ctx := context.Background()
		err := container.HandleConnectFrame(ctx, mockFrame, testCallback)

		if err == nil {
			t.Error("Expected error but got none")
		}
		if container.State != domain.ContainerCreated {
			t.Errorf("Expected state %v, got %v", domain.ContainerCreated, container.State)
		}
		if callbackCount != 0 {
			t.Errorf("Expected 0 callback calls, got %d", callbackCount)
		}
	})
}

func TestContainer_HandleOpenRcvdFrame(t *testing.T) {
	t.Run("HandleOpenRcvdFrame_ValidState_Success", func(t *testing.T) {
		container := NewContainer("container123", "client123")
		container.State = domain.ContainerOpenSent

		payload := mocks.NewMockOpenRcvdFramePayload(t)
		mockFrame := mocks.NewMockFrame(t)
		mockFrame.On("GetPayload").Return(payload)

		err := container.HandleOpenRcvdFrame(mockFrame)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if container.State != domain.ContainerReserved {
			t.Errorf("Expected state %v, got %v", domain.ContainerReserved, container.State)
		}
	})

	t.Run("HandleOpenRcvdFrame_InvalidState_Error", func(t *testing.T) {
		container := NewContainer("container123", "client123")
		container.State = domain.ContainerCreated

		mockFrame := mocks.NewMockFrame(t)

		err := container.HandleOpenRcvdFrame(mockFrame)

		if err == nil {
			t.Error("Expected error but got none")
		}
		if container.State != domain.ContainerCreated {
			t.Errorf("Expected state %v, got %v", domain.ContainerCreated, container.State)
		}
	})
}

func TestContainer_HandleSubscribeFrame(t *testing.T) {
	t.Run("HandleSubscribeFrame_ValidState_Success", func(t *testing.T) {
		container := NewContainer("container123", "client123")
		container.State = domain.ContainerConnected

		topic := "test.topic"
		payload := mocks.NewMockSubscribeFramePayload(t)
		payload.On("GetTopic").Return(topic)

		mockFrame := mocks.NewMockFrame(t)
		mockFrame.On("GetPayload").Return(payload)

		testCallback := func(ctx context.Context, frame domain.Frame, clientID domain.ID) error {
			return nil
		}

		ctx := context.Background()
		err := container.HandleSubscribeFrame(ctx, mockFrame, testCallback)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if _, exists := container.ChannelsByTopic[topic]; !exists {
			t.Errorf("Expected channel for topic %s to be created", topic)
		}
	})

	t.Run("HandleSubscribeFrame_InvalidState_Error", func(t *testing.T) {
		container := NewContainer("container123", "client123")
		container.State = domain.ContainerOpenSent

		mockFrame := mocks.NewMockFrame(t)

		testCallback := func(ctx context.Context, frame domain.Frame, clientID domain.ID) error {
			return nil
		}

		ctx := context.Background()
		err := container.HandleSubscribeFrame(ctx, mockFrame, testCallback)

		if err == nil {
			t.Error("Expected error but got none")
		}
	})
}

func TestContainer_HandleFrame(t *testing.T) {
	t.Run("HandleFrame_UnsupportedType_Error", func(t *testing.T) {
		container := NewContainer("container123", "client123")
		container.State = domain.ContainerConnected

		mockFrame := mocks.NewMockFrame(t)
		mockFrame.On("GetType").Return(domain.FrameTypeMessage)

		testCallback := func(ctx context.Context, frame domain.Frame, clientID domain.ID) error {
			return nil
		}

		ctx := context.Background()
		err := container.HandleFrame(ctx, mockFrame, testCallback)

		if err == nil {
			t.Error("Expected error but got none")
		}
	})
}
