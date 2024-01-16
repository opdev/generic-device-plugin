package plugin

import (
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

func (d *EdgeDevicePlugin) allocateContainerRequests(deviceIDs []string) *pluginapi.ContainerAllocateResponse {
	response := new(pluginapi.ContainerAllocateResponse)

	for _, id := range deviceIDs {
		device := d.rm.GetDeviceByID(id)
		response.Mounts = append(
			response.Mounts,
			&pluginapi.Mount{
				ContainerPath: device.DevPath,
				HostPath:      device.DevPath,
				ReadOnly:      false,
			},
		)

		response.Devices = append(
			response.Devices,
			&pluginapi.DeviceSpec{
				ContainerPath: device.DevPath,
				HostPath:      device.DevPath,
				Permissions:   "rwm",
			},
		)
	}

	return nil
}
