package core

import (
	"context"
	"log/slog"
)

type Broker struct {
	Logger *slog.Logger
}

func (b *Broker) Start(ctx context.Context) error {
	b.Logger.Info("Starting Broker Component")

	return nil
}

func (b *Broker) Stop(ctx context.Context) error {
	b.Logger.Info("Stopping Broker Component")

	return nil
}
