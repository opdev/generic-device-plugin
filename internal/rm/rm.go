package rm

import (
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

type Device struct {
	Name      string
	Major     string
	Minor     string
	BusNum    string
	DevNum    string
	DevName   string
	DevPath   string
	Serial    string
	ProductId string
	VenderId  string
}

type ResourceManager interface {
	Discover() ([]Device, error)
	// AllocateDevice function accepts the arguments:
	// AvailableDeviceIDs ([]string)
	// MustIncludeDeviceIDs ([]string)
	// AllocationSize (int32)
	AllocateDevice(availableDeviceIds, mustIncludeDeviceIDs []string, size int32)
	Devices() []*pluginapi.Device
	GetDeviceByID(id string) *Device
}
