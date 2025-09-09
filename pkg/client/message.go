package client

import (
	"time"

	"github.com/hoppermq/hopper/pkg/domain"
)

type Message struct {
	ID domain.ID
	Topic string
	Content []byte
	Headers map[string]string

	Timestamp time.Time
	SourceID domain.ID

	SubscriptionID domain.ID

	QoS uint8

	protocol string
	originalFrame any
}
