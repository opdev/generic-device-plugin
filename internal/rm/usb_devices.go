package rm

import (
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

var _ ResourceManager = &usbResourceManager{}

type usbResourceManager struct {
	devices []Device
}

// NewManager initializes a Manager which implements
// the resource manager interface
func NewUSBResourceManager() *usbResourceManager {
	devs, _ := findDevices()
	return &usbResourceManager{
		devices: devs,
	}
}

func (r *usbResourceManager) Discover() ([]Device, error) {
	devices, err := findDevices()
	if err != nil {
		return nil, err
	}

	r.devices = devices
	return r.devices, nil
}

func (r *usbResourceManager) Devices() []*pluginapi.Device {
	devices, err := findDevices()
	if err != nil {
		return nil
	}

	var result []*pluginapi.Device
	for _, device := range devices {
		result = append(result, &pluginapi.Device{
			ID:     device.Name,
			Health: pluginapi.Healthy,
		})
	}

	return result
}

func (r *usbResourceManager) GetDeviceByID(id string) *Device {
	devices, err := findDevices()
	if err != nil {
		return nil
	}

	for _, device := range devices {
		if device.Name == id {
			return &device
		}
	}

	return nil
}

// AllocateDevice function accepts the arguments:
// AvailableDeviceIDs ([]string)
// MustIncludeDeviceIDs ([]string)
// AllocationSize (int32)
func (r *usbResourceManager) AllocateDevice(availableDeviceIds []string, mustIncludeDeviceIds []string, allocationSize int32) {
	panic("not implemented") // TODO: Implement
}
