package main

import (
	"context"
	"log/slog"
	"net"

	"github.com/gin-gonic/gin" // should not exist here

	"github.com/hoppermq/hopper/internal/events"
	httpService "github.com/hoppermq/hopper/internal/http"

	"github.com/zixyos/glog"

	"github.com/hoppermq/hopper/internal/application"
	"github.com/hoppermq/hopper/internal/config"
	"github.com/hoppermq/hopper/internal/mq"
	handler "github.com/hoppermq/hopper/internal/mq/transport/tcp"
)

const appName = "Hopper"

func main() {
	ctx := context.Background()
	cfg, err := config.New(appName)
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

	hopperMQService := mq.New(
		mq.WithLogger(logger),
		mq.WithTCP( // should be a more generic transport configuration injection.
			ctx,
			handler.WithListener(conf),
			handler.WithLogger(logger),
		),
	)

	httpEngine := gin.New()

	httpServer := httpService.NewHTTPServer(
		httpService.WithLogger(logger),
		httpService.WithEngine(httpEngine),
	)

	logger.Info("Hey Welcome to HOPPER")
	application.New(
		application.WithConfiguration(cfg),
		application.WithLogger(logger),
		application.WithEventBus(
			events.WithConfig(cfg),
		), // should create it with no parameter
		application.WithService(hopperMQService, httpServer),
	).Start()
}
