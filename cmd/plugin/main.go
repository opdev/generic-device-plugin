package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/opdev/generic-device-plugin/internal/plugin"
)

func main() {
	devplugin := plugin.NewEdgeDevicePlugin()
	ctx := context.Background()

	if err := devplugin.Run(ctx); err != nil {
		log.Fatalf("error starting device plugin server; %+v\n", err)
	}

	log.Println("Device plugin successfully started")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	for {
		select {
		case <-sigCh:
			log.Println("Device plugin is terminating...")
			devplugin.Stop()
			os.Exit(0)
		}
	}
}
