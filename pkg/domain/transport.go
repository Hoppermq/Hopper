package domain

import "context"

type Transport interface {
	Service
	HandleConnection(ctx context.Context) error
	Start(ctx context.Context) error
}
