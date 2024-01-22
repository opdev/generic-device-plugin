package rm

import (
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

var _ ResourceManager = &usbResourceManager{}

type usbResourceManager struct {
	devices map[string]Device
}

// USBResourceManager initializes a new USB resource manager
// which implements the resource manager interface
func USBResourceManager() *usbResourceManager {
	devs, err := findDevices()
	if err != nil {
		return nil
	}

	return &usbResourceManager{
		devices: mapDevices(devs),
	}
}

func (r *usbResourceManager) Devices() []*pluginapi.Device {
	devices, err := findDevices()
	if err != nil {
		return nil
	}

	var result []*pluginapi.Device
	for _, device := range devices {
		result = append(result, device.Device)
	}

	return result
}

func (r *usbResourceManager) DeviceMap() map[string]Device {
	return r.devices
}

func (r *usbResourceManager) GetDeviceByID(id string) (Device, bool) {
	dev, ok := r.devices[id]
	return dev, ok
}

// AllocateDevice function accepts the arguments:
// AvailableDeviceIDs ([]string)
// MustIncludeDeviceIDs ([]string)
// AllocationSize (int32)
func (r *usbResourceManager) AllocateDevice(availableDeviceIds []string, mustIncludeDeviceIds []string, allocationSize int32) {
	panic("not implemented") // TODO: Implement
}
