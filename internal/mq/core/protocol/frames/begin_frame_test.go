package frames

import (
	"testing"

	"github.com/hoppermq/hopper/pkg/domain"
)

func TestCreateBeginFramePayload(t *testing.T) {
	tests := []struct {
		name           string
		sourceID       domain.ID
		containerID    domain.ID
		remoteChannel  uint16
		nextOutgoingID uint32
		incomingWindow uint32
		outgoingWindow uint32
		wantErr        bool
		validate       func(t *testing.T, payload *BeginFramePayload, err error)
	}{
		{
			name:           "CreateBeginFramePayload_ValidParameters",
			sourceID:       "client123",
			containerID:    "container456",
			remoteChannel:  0,
			nextOutgoingID: 0,
			incomingWindow: 1000,
			outgoingWindow: 1000,
			wantErr:        false,
			validate: func(t *testing.T, payload *BeginFramePayload, err error) {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}
				if payload.GetSourceID() != "client123" {
					t.Errorf("Expected source ID 'client123', got %v", payload.GetSourceID())
				}
				if payload.GetContainerID() != "container456" {
					t.Errorf("Expected container ID 'container456', got %v", payload.GetContainerID())
				}
				if payload.GetRemoteChannel() != 0 {
					t.Errorf("Expected remote channel 0, got %v", payload.GetRemoteChannel())
				}
				if payload.GetNextOutgoingID() != 0 {
					t.Errorf("Expected next outgoing ID 0, got %v", payload.GetNextOutgoingID())
				}
				if payload.GetIncomingWindow() != 1000 {
					t.Errorf("Expected incoming window 1000, got %v", payload.GetIncomingWindow())
				}
				if payload.GetOutgoingWindow() != 1000 {
					t.Errorf("Expected outgoing window 1000, got %v", payload.GetOutgoingWindow())
				}
			},
		},
		{
			name:           "CreateBeginFramePayload_DifferentValues",
			sourceID:       "client789",
			containerID:    "container012",
			remoteChannel:  5,
			nextOutgoingID: 100,
			incomingWindow: 2000,
			outgoingWindow: 1500,
			wantErr:        false,
			validate: func(t *testing.T, payload *BeginFramePayload, err error) {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}
				if payload.GetSourceID() != "client789" {
					t.Errorf("Expected source ID 'client789', got %v", payload.GetSourceID())
				}
				if payload.GetContainerID() != "container012" {
					t.Errorf("Expected container ID 'container012', got %v", payload.GetContainerID())
				}
				if payload.GetRemoteChannel() != 5 {
					t.Errorf("Expected remote channel 5, got %v", payload.GetRemoteChannel())
				}
				if payload.GetNextOutgoingID() != 100 {
					t.Errorf("Expected next outgoing ID 100, got %v", payload.GetNextOutgoingID())
				}
				if payload.GetIncomingWindow() != 2000 {
					t.Errorf("Expected incoming window 2000, got %v", payload.GetIncomingWindow())
				}
				if payload.GetOutgoingWindow() != 1500 {
					t.Errorf("Expected outgoing window 1500, got %v", payload.GetOutgoingWindow())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := CreateBeginFramePayload(
				&PayloadHeader{},
				tt.sourceID,
				tt.containerID,
				tt.remoteChannel,
				tt.nextOutgoingID,
				tt.incomingWindow,
				tt.outgoingWindow,
			)

			tt.validate(t, payload, nil)
		})
	}
}

func TestBeginFramePayload_Sizer(t *testing.T) {
	tests := []struct {
		name        string
		sourceID    domain.ID
		containerID domain.ID
		wantErr     bool
		validate    func(t *testing.T, size uint16, err error)
	}{
		{
			name:        "Sizer_ShortIDs",
			sourceID:    "c1",
			containerID: "c2",
			wantErr:     false,
			validate: func(t *testing.T, size uint16, err error) {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}
				expectedSize := uint16(2 + 2 + 2 + 4 + 4 + 4 + 2)
				if size != expectedSize {
					t.Errorf("Expected size %v, got %v", expectedSize, size)
				}
			},
		},
		{
			name:        "Sizer_LongerIDs",
			sourceID:    "client123456",
			containerID: "container789012",
			wantErr:     false,
			validate: func(t *testing.T, size uint16, err error) {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}
				expectedSize := uint16(12 + 15 + 2 + 4 + 4 + 4 + 2)
				if size != expectedSize {
					t.Errorf("Expected size %v, got %v", expectedSize, size)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := CreateBeginFramePayload(
				&PayloadHeader{},
				tt.sourceID,
				tt.containerID,
				0, 0, 1000, 1000,
			)
			size := payload.Sizer()

			tt.validate(t, size, nil)
		})
	}
}

func TestCreateBeginFrame(t *testing.T) {
	tests := []struct {
		name           string
		doff           domain.DOFF
		sourceID       domain.ID
		containerID    domain.ID
		remoteChannel  uint16
		nextOutgoingID uint32
		incomingWindow uint32
		outgoingWindow uint32
		wantErr        bool
		validate       func(t *testing.T, frame *Frame, err error)
	}{
		{
			name:           "CreateBeginFrame_ValidParameters",
			doff:           domain.DOFF4,
			sourceID:       "client123",
			containerID:    "container456",
			remoteChannel:  0,
			nextOutgoingID: 0,
			incomingWindow: 1000,
			outgoingWindow: 1000,
			wantErr:        false,
			validate: func(t *testing.T, frame *Frame, err error) {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}
				if frame == nil {
					t.Fatal("Expected frame but got nil")
					return
				}
				if frame.GetType() != domain.FrameTypeBegin {
					t.Errorf("Expected frame type %v, got %v", domain.FrameTypeBegin, frame.GetType())
				}
				beginPayload, ok := frame.GetPayload().(domain.BeginFramePayload)
				if !ok {
					t.Error("Expected BeginFramePayload")
					return
				}
				if beginPayload.GetSourceID() != "client123" {
					t.Errorf("Expected source ID 'client123', got %v", beginPayload.GetSourceID())
				}
				if beginPayload.GetContainerID() != "container456" {
					t.Errorf("Expected container ID 'container456', got %v", beginPayload.GetContainerID())
				}
			},
		},
		{
			name:           "CreateBeginFrame_DifferentDOFF",
			doff:           domain.DOFF3,
			sourceID:       "client789",
			containerID:    "container012",
			remoteChannel:  5,
			nextOutgoingID: 100,
			incomingWindow: 2000,
			outgoingWindow: 1500,
			wantErr:        false,
			validate: func(t *testing.T, frame *Frame, err error) {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}
				if frame == nil {
					t.Fatal("Expected frame but got nil")
					return
				}
				if frame.GetType() != domain.FrameTypeBegin {
					t.Errorf("Expected frame type %v, got %v", domain.FrameTypeBegin, frame.GetType())
				}
				beginPayload, ok := frame.GetPayload().(domain.BeginFramePayload)
				if !ok {
					t.Error("Expected BeginFramePayload")
					return
				}
				if beginPayload.GetSourceID() != "client789" {
					t.Errorf("Expected source ID 'client789', got %v", beginPayload.GetSourceID())
				}
				if beginPayload.GetContainerID() != "container012" {
					t.Errorf("Expected container ID 'container012', got %v", beginPayload.GetContainerID())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frame, err := CreateBeginFrame(
				tt.doff,
				tt.sourceID,
				tt.containerID,
				tt.remoteChannel,
				tt.nextOutgoingID,
				tt.incomingWindow,
				tt.outgoingWindow,
			)

			tt.validate(t, frame, err)
		})
	}
}

func TestBeginFrameIntegration(t *testing.T) {
	tests := []struct {
		name     string
		wantErr  bool
		validate func(t *testing.T, frame *Frame, err error)
	}{
		{
			name:    "BeginFrame_FullIntegration",
			wantErr: false,
			validate: func(t *testing.T, frame *Frame, err error) {
				sourceID := domain.ID("integration-client")
				containerID := domain.ID("integration-container")
				remoteChannel := uint16(10)
				nextOutgoingID := uint32(500)
				incomingWindow := uint32(3000)
				outgoingWindow := uint32(2500)

				if err != nil {
					t.Fatalf("Failed to create Begin frame: %v", err)
				}

				if frame.GetType() != domain.FrameTypeBegin {
					t.Errorf("Expected Begin frame type, got %v", frame.GetType())
				}

				beginPayload, ok := frame.GetPayload().(domain.BeginFramePayload)
				if !ok {
					t.Fatal("Expected BeginFramePayload")
				}

				if beginPayload.GetSourceID() != sourceID {
					t.Errorf("Source ID mismatch: expected %v, got %v", sourceID, beginPayload.GetSourceID())
				}

				if beginPayload.GetContainerID() != containerID {
					t.Errorf("Container ID mismatch: expected %v, got %v", containerID, beginPayload.GetContainerID())
				}

				if beginPayload.GetRemoteChannel() != remoteChannel {
					t.Errorf("Remote channel mismatch: expected %v, got %v", remoteChannel, beginPayload.GetRemoteChannel())
				}

				if beginPayload.GetNextOutgoingID() != nextOutgoingID {
					t.Errorf("Next outgoing ID mismatch: expected %v, got %v", nextOutgoingID, beginPayload.GetNextOutgoingID())
				}

				if beginPayload.GetIncomingWindow() != incomingWindow {
					t.Errorf("Incoming window mismatch: expected %v, got %v", incomingWindow, beginPayload.GetIncomingWindow())
				}

				if beginPayload.GetOutgoingWindow() != outgoingWindow {
					t.Errorf("Outgoing window mismatch: expected %v, got %v", outgoingWindow, beginPayload.GetOutgoingWindow())
				}

				if !frame.CanHandle(domain.FrameTypeBegin) {
					t.Error("Frame should be able to handle Begin frame type")
				}

				expectedSize := beginPayload.Sizer()
				if expectedSize == 0 {
					t.Error("Expected non-zero payload size")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frame, err := CreateBeginFrame(
				domain.DOFF4,
				"integration-client",
				"integration-container",
				uint16(10),
				uint32(500),
				uint32(3000),
				uint32(2500),
			)

			tt.validate(t, frame, err)
		})
	}
}
