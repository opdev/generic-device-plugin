package usb

import (
	"github.com/OchiengEd/edge-device-plugin/internal/rm"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

var _ rm.ResourceManager = &Manager{}

type Manager struct {
	devices []rm.Device
}

// NewManager initializes a Manager which implements
// the resource manager interface
func NewManager() *Manager {
	devs, _ := findDevices()
	return &Manager{
		devices: devs,
	}
}

func (r *Manager) Discover() ([]rm.Device, error) {
	devices, err := findDevices()
	if err != nil {
		return nil, err
	}

	r.devices = devices
	return r.devices, nil
}

func (r *Manager) Devices() []*pluginapi.Device {
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

func (r *Manager) GetDeviceByID(id string) *rm.Device {
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
func (r *Manager) AllocateDevice(availableDeviceIds []string, mustIncludeDeviceIds []string, allocationSize int32) {
	panic("not implemented") // TODO: Implement
}
