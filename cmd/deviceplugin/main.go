package main

import (
	"context"
	"fmt"

	"github.com/OchiengEd/edge-device-plugin/intenral/plugin"
)

func main() {
	fmt.Println("vim-go")
	devplugin := plugin.NewEdgeDevicePlugin()

	ctx := context.Background()
	if err := devplugin.Start(ctx); err != nil {
		fmt.Printf("error starting device plugin server; %+v\n", err)
	}

	fmt.Println("Successfully registered")

	select {
	case msg := <-ctx.Done():
		fmt.Printf("Context cancellation; %+v\n", msg)
	}
}
