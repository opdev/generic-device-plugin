package plugin

import (
	"context"
	"log"

	"github.com/fsnotify/fsnotify"
	"google.golang.org/grpc"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

func (d *EdgeDevicePlugin) Run(ctx context.Context) error {
	// Start device plugin server and register with Kubelet
	if err := d.Start(ctx); err != nil {
		log.Fatalf("error starting device plugin server; %+v\n", err)
	}

	// Setup monitoring for device plugin unix socket
	return d.monitorSocket(ctx)
}

func (d *EdgeDevicePlugin) monitorSocket(ctx context.Context) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Has(fsnotify.Remove) && event.Name == d.socket {
					d.Stop()
					d.quit <- true

					d.grpcServer = grpc.NewServer()
					if err := d.Start(ctx); err != nil {
						log.Fatalf("error restarting device plugin...")
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				if err != nil {
					log.Printf("fsnotify encountered an error: %+v\n", err)
				}
			}
		}
	}()

	return watcher.Add(pluginapi.DevicePluginPath)
}
