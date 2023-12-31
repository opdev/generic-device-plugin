package plugin

import (
	"context"
	"log"
	"net"
	"os"
	"path"
	"time"

	"google.golang.org/grpc"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

const (
	LitmusSocketName = "litmus.sock"
)

type EdgeDevicePlugin struct {
	socket     string
	grpcServer *grpc.Server
	Heartbeat  chan string
}

func NewEdgeDevicePlugin() *EdgeDevicePlugin {
	return &EdgeDevicePlugin{
		// TODO: remove temporary socket
		// socket:     path.Join(pluginapi.DevicePluginPath, LitmusSocketName),
		socket:     path.Join("/Users/edmund.ochieng", LitmusSocketName),
		grpcServer: grpc.NewServer(),
		Heartbeat:  make(chan string),
	}
}

func (d *EdgeDevicePlugin) Start(ctx context.Context) error {
	if err := d.Serve(ctx); err != nil {
		return err
	}

	err := d.RegisterWithKubelet(ctx, 5*time.Second)
	if err != nil {
		log.Printf("error registering with kubelet; %+v\n", err)
		return err
	}

	return nil
}

func (d *EdgeDevicePlugin) Stop() {
	if d != nil && d.grpcServer != nil {
		d.grpcServer.Stop()
	}
	os.Remove(d.socket)
	close(d.Heartbeat)
}

func (d *EdgeDevicePlugin) Serve(ctx context.Context) error {
	os.Remove(d.socket)

	sock, err := net.Listen("unix", d.socket)
	if err != nil {
		return err
	}
	pluginapi.RegisterDevicePluginServer(d.grpcServer, d)

	go func() {
		if err := d.grpcServer.Serve(sock); err == nil {
			log.Println("gRPC server crashed while starting...")
		}
	}()

	conn, err := dialConn(ctx, d.socket, 5*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Println("gRPC server successfully started")

	return nil
}

// GetDevicePluginOptions returns options to be communicated with Device
// Manager
func (d *EdgeDevicePlugin) GetDevicePluginOptions(_ context.Context, _ *pluginapi.Empty) (*pluginapi.DevicePluginOptions, error) {
	return &pluginapi.DevicePluginOptions{
		PreStartRequired:                false,
		GetPreferredAllocationAvailable: false,
	}, nil
}

var _ pluginapi.DevicePluginServer = &EdgeDevicePlugin{}

// ListAndWatch returns a stream of List of Devices
// Whenever a Device state change or a Device disappears, ListAndWatch
// returns the new list
func (d *EdgeDevicePlugin) ListAndWatch(_ *pluginapi.Empty, srv pluginapi.DevicePlugin_ListAndWatchServer) error {
	if err := srv.Send(&pluginapi.ListAndWatchResponse{Devices: []*pluginapi.Device{}}); err != nil {
		return err
	}

	for {
		select {
		case <-d.Heartbeat:
			if err := srv.Send(&pluginapi.ListAndWatchResponse{Devices: []*pluginapi.Device{}}); err != nil {
				return nil
			}
		}
	}

	return nil
}

// GetPreferredAllocation returns a preferred set of devices to allocate
// from a list of available ones. The resulting preferred allocation is not
// guaranteed to be the allocation ultimately performed by the
// devicemanager. It is only designed to help the devicemanager make a more
// informed allocation decision when possible.
func (d *EdgeDevicePlugin) GetPreferredAllocation(_ context.Context, prefs *pluginapi.PreferredAllocationRequest) (*pluginapi.PreferredAllocationResponse, error) {
	for _, req := range prefs.ContainerRequests {
		log.Printf("Preferred requests: %+v\n", req)
	}
	return &pluginapi.PreferredAllocationResponse{}, nil
}

// Allocate is called during container creation so that the Device
// Plugin can run device specific operations and instruct Kubelet
// of the steps to make the Device available in the container
func (d *EdgeDevicePlugin) Allocate(_ context.Context, reqs *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error) {

	for _, req := range reqs.ContainerRequests {
		_ = req.GetDevicesIDs()
		log.Printf("Request received: %+v\n", req)
	}

	return &pluginapi.AllocateResponse{}, nil
}

// PreStartContainer is called, if indicated by Device Plugin during registeration phase,
// before each container start. Device plugin can run device specific operations
// such as resetting the device before making devices available to the container
func (d *EdgeDevicePlugin) PreStartContainer(_ context.Context, prereq *pluginapi.PreStartContainerRequest) (*pluginapi.PreStartContainerResponse, error) {
	return &pluginapi.PreStartContainerResponse{}, nil
}
