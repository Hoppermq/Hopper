package core

import "log/slog"

type Broker struct {
	Logger *slog.Logger
}

func (b *Broker) Start() {
	b.Logger.Info("Broker Started")
}
