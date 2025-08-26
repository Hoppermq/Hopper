package domain

import "errors"

var (
	ErrInvalidHeader        = errors.New("invalid header provided")
	ErrInvalidPayload       = errors.New("invalid payload provided")
	ErrUnsupportedFrameType = errors.New("unsupported frame type")
)
