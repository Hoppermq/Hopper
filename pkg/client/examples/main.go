package main

import (
	"context"

	"github.com/hoppermq/hopper/pkg/client"
	"github.com/zixyos/glog"
)

func main() {
	logger, err := glog.NewDefault()
	if err != nil {
		return
	}

	sdk := client.NewClient(
		client.WithConfig(nil),
		client.WithLogger(logger),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err = sdk.Run(ctx); err != nil {
		return
	}

	if err = sdk.Stop(ctx); err != nil {
		return
	}
}
