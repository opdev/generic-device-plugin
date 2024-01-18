package main

import (
	"context"

	"github.com/OchiengEd/edge-device-plugin/internal/plugin"
	"github.com/golang/glog"
)

func main() {
	devplugin := plugin.NewEdgeDevicePlugin()

	ctx := context.Background()
	if err := devplugin.Start(ctx); err != nil {
		glog.Error("error starting device plugin server; %+v\n", err)
	}

	glog.Info("Device plugin successfully started")

	msg := <-ctx.Done()
	glog.Info("Context cancellation; %+v\n", msg)
}
