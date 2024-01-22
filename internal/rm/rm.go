package rm

import (
	"fmt"

	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

type ResourceManager interface {
	Devices() []*pluginapi.Device
	DeviceMap() map[string]Device
	GetDeviceByID(id string) (Device, bool)
	// AllocateDevice function accepts the arguments:
	// AvailableDeviceIDs ([]string)
	// MustIncludeDeviceIDs ([]string)
	// AllocationSize (int32)
	AllocateDevice(available, required []string, size int32)
}

type Device struct {
	*pluginapi.Device
	Name         string
	Major        string
	Minor        string
	BusNum       string
	DevNum       string
	DevName      string
	DevPath      string
	Serial       string
	ProductId    string
	VendorId     string
	Product      string
	Manufacturer string
}

func (d *Device) UniqueID() string {
	return fmt.Sprintf("%s:%s", d.VendorId, d.ProductId)
}

func (d *Device) Mounts() []*pluginapi.Mount {
	return []*pluginapi.Mount{
		{
			HostPath:      d.DevPath,
			ContainerPath: d.DevPath,
			ReadOnly:      false,
		},
	}
}

func (d *Device) DeviceSpec() []*pluginapi.DeviceSpec {
	return []*pluginapi.DeviceSpec{
		{
			HostPath:      d.DevPath,
			ContainerPath: d.DevPath,
			Permissions:   "rwm",
		},
	}
}

func (d *Device) Is(d2 *Device) bool {
	if d.ID != d2.ID {
		return false
	}

	// Compare plugin device
	if d.Device != nil && d2.Device != nil {
		if d.Device.ID != d2.Device.ID {
			return false
		}

		if d.Device.Health != d2.Device.Health {
			return false
		}
	}

	if d.DevPath != d2.DevPath {
		return false
	}

	if d.UniqueID() != d2.UniqueID() {
		return false
	}

	if d.Serial != d2.Serial {
		return false
	}

	if d.Manufacturer != d2.Manufacturer {
		return false
	}

	if d.Product != d2.Product {
		return false
	}

	return true
}
