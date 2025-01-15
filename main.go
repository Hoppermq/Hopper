package main

import (
	"log/slog"

	"github.com/hoppermq/hopper/internal/application"
	"github.com/hoppermq/hopper/internal/config"
	"github.com/zixyos/glog"
)

const appName = "Hopper";

func main() {
  _, err := config.New(appName);
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
  );

  if err != nil {
    panic(err)
  }

  logger.Info("Hey Welcome to HOPPER");
  application.New(
    application.WithLogger(logger),
  ).Start();
}
