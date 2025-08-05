package main

import (
	"context"
	"log/slog"
	"net"

	"github.com/hoppermq/hopper/internal/application"
	"github.com/hoppermq/hopper/internal/config"
	"github.com/hoppermq/hopper/internal/handler"
	"github.com/hoppermq/hopper/internal/mq"
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

	lconf := &net.ListenConfig{}

	tcpHandler, err := handler.NewTCP(
		handler.WithContext(ctx),
		handler.WithListener(*lconf),
		handler.WithLogger(*&logger),
	)

	if err != nil {
		panic(err)
	}

	mq := mq.New(
		mq.WithLogger(logger),
		mq.WithTCP(tcpHandler),
	)

	logger.Info("Hey Welcome to HOPPER")
	application.New(
		application.WithLogger(logger),
		application.WithService(mq),
	).Start()
}
