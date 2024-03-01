package plugin

import (
	"context"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

func dialConn(ctx context.Context, socket string, timeout time.Duration) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var d net.Dialer
	conn, err := grpc.DialContext(ctx, socket,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return d.DialContext(ctx, "unix", s)
		}),
	)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (d *EdgeDevicePlugin) RegisterWithKubelet(ctx context.Context, timeout time.Duration) error {
	conn, err := dialConn(ctx, pluginapi.KubeletSocket, 5*time.Second)
	if err != nil {
		return err
	}

	// Get device plugin options
	options, err := d.GetDevicePluginOptions(ctx, nil)
	if err != nil {
		return err
	}

	client := pluginapi.NewRegistrationClient(conn)
	_, err = client.Register(ctx,
		&pluginapi.RegisterRequest{
			Version:      pluginapi.Version,
			ResourceName: "example.io/device",
			Endpoint:     GenericSocketName,
			Options:      options,
		},
	)

	return err
}
