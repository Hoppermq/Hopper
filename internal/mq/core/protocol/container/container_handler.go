package container

import (
	"context"
	"fmt"

	"github.com/hoppermq/hopper/internal/common"
	"github.com/hoppermq/hopper/internal/mq/core/protocol/frames"
	"github.com/hoppermq/hopper/pkg/domain"
)

// FrameSendCallback represents a callback function for sending frames back to clients
type FrameSendCallback func(ctx context.Context, frame domain.Frame, clientID domain.ID) error

// HandleConnectFrame handles Connect frame and creates Begin frame response using callback approach
func (ctr *Container) HandleConnectFrame(ctx context.Context, frame domain.Frame, sendCallback FrameSendCallback) error {
	if ctr.State != domain.ContainerOpenSent {
		return fmt.Errorf("invalid container state for Connect frame: expected %s, got %s",
			domain.ContainerOpenSent, ctr.State)
	}

	connectPayload, ok := frame.GetPayload().(domain.ConnectFramePayload)
	if !ok {
		return fmt.Errorf("invalid payload type for Connect frame")
	}

	beginFrame, err := ctr.createBeginFrame(connectPayload.GetSourceID())
	if err != nil {
		return fmt.Errorf("failed to create Begin frame: %w", err)
	}

	_ = ctr.CreateChannel("__temp__", common.GenerateIdentifier)

	ctr.State = domain.ContainerConnected

	return sendCallback(ctx, beginFrame, connectPayload.GetSourceID())
}

func (ctr *Container) HandleOpenRcvdFrame(frame domain.Frame) error {
	if ctr.State != domain.ContainerOpenSent {
		return fmt.Errorf("invalid container state for OpenRcvd frame: expected %s, got %s",
			domain.ContainerOpenSent, ctr.State)
	}

	_, ok := frame.GetPayload().(domain.OpenRcvdFramePayload)
	if !ok {
		return fmt.Errorf("invalid payload type for OpenRcvd frame")
	}

	ctr.State = domain.ContainerReserved

	return nil
}

// HandleSubscribeFrame handles Subscribe frame and creates channels for topic subscription
func (ctr *Container) HandleSubscribeFrame(ctx context.Context, frame domain.Frame, sendCallback FrameSendCallback) error {
	if ctr.State != domain.ContainerConnected {
		return fmt.Errorf("invalid container state for Subscribe frame: expected %s, got %s",
			domain.ContainerConnected, ctr.State)
	}

	subscribePayload, ok := frame.GetPayload().(domain.SubscribeFramePayload)
	if !ok {
		return fmt.Errorf("invalid payload type for Subscribe frame")
	}

	topic := subscribePayload.GetTopic()
	if _, exists := ctr.ChannelsByTopic[topic]; !exists {
		ctr.CreateChannel(topic, common.GenerateIdentifier)
	}

	return nil
}

// createBeginFrame creates a Begin frame for this container
func (ctr *Container) createBeginFrame(sourceID domain.ID) (domain.Frame, error) {
	beginFrame, err := frames.CreateBeginFrame(
		domain.DOFF4,
		sourceID,
		ctr.ID,
		0,
		0,
		1000,
		1000,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Begin frame: %w", err)
	}

	return beginFrame, nil
}

func (ctr *Container) HandleFrame(ctx context.Context, frame domain.Frame, sendCallback FrameSendCallback) error {
	frameType := frame.GetType()

	switch frameType {
	case domain.FrameTypeOpenRcvd:
		return ctr.HandleOpenRcvdFrame(frame)
	case domain.FrameTypeConnect:
		return ctr.HandleConnectFrame(ctx, frame, sendCallback)
	case domain.FrameTypeSubscribe:
		return ctr.HandleSubscribeFrame(ctx, frame, sendCallback)
	default:
		return fmt.Errorf("unsupported frame type: %v", frameType)
	}
}

func (ctr *Container) GetClientID() domain.ID {
	return ctr.ClientID
}
