package main

import (
	"context"
	"log/slog"
	"net"
	"os"

	"github.com/gin-gonic/gin" // should not exist here
	"github.com/hoppermq/hopper/internal/ui"

	"github.com/hoppermq/hopper/internal/events"
	httpService "github.com/hoppermq/hopper/internal/http"

	"github.com/zixyos/glog"

	"github.com/hoppermq/hopper/internal/application"
	"github.com/hoppermq/hopper/internal/config"
	"github.com/hoppermq/hopper/internal/mq"
	"github.com/hoppermq/hopper/internal/mq/core"
	handler "github.com/hoppermq/hopper/internal/mq/transport/tcp"
)

const (
	appName       = "Hopper"
	maxBufferSize = 1024
)

func main() {
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
		slog.New(slog.NewJSONHandler(os.Stdout, nil)).Error("Failed to create logger", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()
	cfg, err := config.New(appName)
	if err != nil {
		logger.Warn("failed to load config", "error", err)
	}

	logger.Info("welcome to " + appName)

	eventBus := events.NewEventBus(maxBufferSize)
	// TODO : HOP-000 should use app config directly
	conf := &net.ListenConfig{}
	tcpTransport, err := handler.NewTCP(
		ctx,
		handler.WithListener(conf),
		handler.WithLogger(logger),
	)
	if err != nil {
		logger.Error("failed to create transport", "error", err)
		os.Exit(1)
	}

	// TODO: HOP-000 should use composite pattern func.
	broker := core.NewBroker(
		logger,
		eventBus,
		tcpTransport,
	)

	hopperMQService := mq.New(
		mq.WithLogger(logger),
		mq.WithBroker(broker),
		mq.WithTransport(tcpTransport),
		mq.WithEventBus(eventBus),
	)

	httpEngine := gin.New()
	httpServer := httpService.NewHTTPServer(
		httpService.WithLogger(logger),
		httpService.WithEngine(httpEngine),
	)

	uiEngine := gin.New()
	uiService := ui.NewHTTPServer(
		ui.WithLogger(logger),
		ui.WithEngine(uiEngine),
	)

	app := application.New(
		application.WithConfiguration(cfg),
		application.WithLogger(logger),
		application.WithEventBus(eventBus),
		application.WithService(hopperMQService, httpServer, uiService),
	)

	app.Start()
}
