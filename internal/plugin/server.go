package plugin

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"path"
	"time"

	"github.com/opdev/generic-device-plugin/internal/rm"
	"google.golang.org/grpc"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

const (
	GenericSocketName = "generic.sock"
	heartbeat         = 10 * time.Second
)

type EdgeDevicePlugin struct {
	rm         rm.ResourceManager
	socket     string
	grpcServer *grpc.Server
	quit       chan interface{}
}

func NewEdgeDevicePlugin() *EdgeDevicePlugin {
	return &EdgeDevicePlugin{
		rm: rm.USBResourceManager(),
		socket: path.Join(
			pluginapi.DevicePluginPath,
			GenericSocketName,
		),
		grpcServer: grpc.NewServer(),
		quit:       make(chan interface{}),
	}
}

func (d *EdgeDevicePlugin) Start(ctx context.Context) error {
	// Start service requests for Device Plugin (gRPC) server
	if err := d.Serve(ctx); err != nil {
		return err
	}

	// Register device plugin with kubelet
	err := d.RegisterWithKubelet(ctx, 5*time.Second)
	if err != nil {
		return err
	}

	return nil
}

func (d *EdgeDevicePlugin) Stop() {
	if d != nil && d.grpcServer != nil {
		d.grpcServer.Stop()
	}
	d.grpcServer = nil
	os.Remove(d.socket)
}

func (d *EdgeDevicePlugin) Serve(ctx context.Context) error {
	_ = os.Remove(d.socket)

	sock, err := net.Listen("unix", d.socket)
	if err != nil {
		return err
	}
	pluginapi.RegisterDevicePluginServer(d.grpcServer, d)

	go func() {
		for {
			select {
			case <-d.quit:
				// Terminate goroutine when a message is written to quit channel
				return
			default:
				if err := d.grpcServer.Serve(sock); err != nil {
					log.Println("error: DevicePlugin server crashed while starting...", err.Error())
				}
			}
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

var _ pluginapi.DevicePluginServer = &EdgeDevicePlugin{}

// GetDevicePluginOptions returns options to be communicated with Device
// Manager
func (d *EdgeDevicePlugin) GetDevicePluginOptions(_ context.Context, _ *pluginapi.Empty) (*pluginapi.DevicePluginOptions, error) {
	return &pluginapi.DevicePluginOptions{
		PreStartRequired:                false,
		GetPreferredAllocationAvailable: false,
	}, nil
}

// ListAndWatch returns a stream of List of Devices
// Whenever a Device state change or a Device disappears, ListAndWatch
// returns the new list
func (d *EdgeDevicePlugin) ListAndWatch(_ *pluginapi.Empty, srv pluginapi.DevicePlugin_ListAndWatchServer) error {
	if err := srv.Send(&pluginapi.ListAndWatchResponse{Devices: d.rm.Devices()}); err != nil {
		return err
	}

	for {
		select {
		case <-time.After(heartbeat):
			if err := srv.Send(&pluginapi.ListAndWatchResponse{Devices: d.rm.Devices()}); err != nil {
				return err
			}
		}
	}
}

// GetPreferredAllocation returns a preferred set of devices to allocate
// from a list of available ones. The resulting preferred allocation is not
// guaranteed to be the allocation ultimately performed by the
// devicemanager. It is only designed to help the devicemanager make a more
// informed allocation decision when possible.
func (d *EdgeDevicePlugin) GetPreferredAllocation(_ context.Context, prefs *pluginapi.PreferredAllocationRequest) (*pluginapi.PreferredAllocationResponse, error) {
	for _, req := range prefs.ContainerRequests {
		fmt.Println(req.AvailableDeviceIDs, req.MustIncludeDeviceIDs, req.AllocationSize)
	}
	return &pluginapi.PreferredAllocationResponse{}, nil
}

// Allocate is called during container creation so that the Device
// Plugin can run device specific operations and instruct Kubelet
// of the steps to make the Device available in the container
func (d *EdgeDevicePlugin) Allocate(_ context.Context, req *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error) {
	var response = &pluginapi.AllocateResponse{
		ContainerResponses: make([]*pluginapi.ContainerAllocateResponse, 0, len(req.ContainerRequests)),
	}

	for _, r := range req.ContainerRequests {
		resp := new(pluginapi.ContainerAllocateResponse)
		// Add devices in request to the response
		for _, id := range r.DevicesIDs {
			dev, ok := d.rm.GetDeviceByID(id)
			if !ok {
				return nil, fmt.Errorf("requested device %q does not exist", id)
			}
			if dev.Device.Health != pluginapi.Healthy {
				return nil, fmt.Errorf("requested device %q is not healthy", id)
			}
			resp.Devices = append(resp.Devices, dev.DeviceSpec()...)
			resp.Mounts = append(resp.Mounts, dev.Mounts()...)
		}
		response.ContainerResponses = append(response.ContainerResponses, resp)
	}

	return response, nil
}

// PreStartContainer is called, if indicated by Device Plugin during registeration phase,
// before each container start. Device plugin can run device specific operations
// such as resetting the device before making devices available to the container
func (d *EdgeDevicePlugin) PreStartContainer(_ context.Context, prereq *pluginapi.PreStartContainerRequest) (*pluginapi.PreStartContainerResponse, error) {
	return &pluginapi.PreStartContainerResponse{}, nil
}
