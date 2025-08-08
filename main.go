package main

import (
	"context"
	"github.com/hoppermq/hopper/internal/events"
	"log/slog"
	"net"

	"github.com/hoppermq/hopper/internal/application"
	"github.com/hoppermq/hopper/internal/config"
	"github.com/hoppermq/hopper/internal/mq"
	handler "github.com/hoppermq/hopper/internal/mq/transport"
	"github.com/zixyos/glog"
)

const appName = "Hopper"

func main() {
	ctx := context.Background()
	_, err := config.New(appName)
	if err != nil {
		panic(err)
	}

	logger, err := glog.New(
		glog.WithLevel(slog.LevelDebug),
		glog.WithFormat(glog.TextFormatter),
		glog.WithTimeStamp(),
		glog.WithReportCaller(),
		glog.WithStyle(
			glog.WithErrorStyle(),
		),
	)
	if err != nil {
		panic(err)
	}

	conf := &net.ListenConfig{}

	eb := events.NewEventBus(1000) // should load from config

	tcpHandler, err := handler.NewTCP(
		handler.WithContext(ctx),
		handler.WithListener(*conf),
		handler.WithLogger(logger),
	)

	if err != nil {
		panic(err)
	}

	hopperMQService := mq.New(
		mq.WithLogger(logger),
		mq.WithTCP(tcpHandler),
	)

	logger.Info("Hey Welcome to HOPPER")
	application.New(
		application.WithLogger(logger),
		application.WithEventBus(eb),
		application.WithService(hopperMQService),
	).Start()
}
