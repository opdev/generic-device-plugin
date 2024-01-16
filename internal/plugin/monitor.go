package plugin

import (
	"log"

	"github.com/fsnotify/fsnotify"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

func monitorKubeletSocket(ch chan<- interface{}) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go func() {
		log.Println("monitoring kubelet socket...")
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Has(fsnotify.Remove) && event.Name == pluginapi.KubeletSocket {
					log.Println("the kubelet socket has been stopped")
				}

				if event.Has(fsnotify.Create) && event.Name == pluginapi.KubeletSocket {
					ch <- true
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
