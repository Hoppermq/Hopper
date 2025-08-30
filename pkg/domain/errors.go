package domain

import "errors"

var (
	// ErrInvalidHeader represent the invalid header provided.
	ErrInvalidHeader = errors.New("invalid header provided")

	// ErrInvalidPayload represent the invalid payload provided.
	ErrInvalidPayload = errors.New("invalid payload provided")

	// ErrUnsupportedFrameType represent the invalid FrameType provided.
	ErrUnsupportedFrameType = errors.New("unsupported frame type")

	// ErrNoServiceAvailable represent the error type when a service is not loaded.
	ErrNoServiceAvailable = errors.New("no service available")
)
