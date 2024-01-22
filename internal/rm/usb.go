package rm

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

func findDevices() ([]Device, error) {
	uevents, err := filepath.Glob("/sys/bus/usb/devices/[0-9]-[0-9]/uevent")
	if err != nil {
		return nil, err
	}

	var devices []Device
	for _, device := range uevents {
		dev := parseDeviceInfo(device)
		devices = append(devices, *dev)
	}

	return devices, nil
}

func mapDevices(devices []Device) map[string]Device {
	var deviceMap = make(map[string]Device)
	for _, dev := range devices {
		deviceMap[dev.DevName] = dev
	}

	return deviceMap
}

func parseDeviceInfo(uevent string) *Device {
	device, err := createDeviceFromUevent(uevent)
	if err != nil {
		return nil
	}

	// Set properties not in the uevent file
	device.Serial = getSerial(device)
	device.ProductId = getProductId(device)
	device.VendorId = getVendorId(device)
	device.Manufacturer = getManufacturer(device)
	device.Product = getProduct(device)

	device.Device = &pluginapi.Device{
		ID:     device.Name,
		Health: pluginapi.Healthy,
	}

	return device
}

func getDeviceNameFromUevent(path string) string {
	re := regexp.MustCompile("/sys/bus/usb/devices/(?P<dev>[0-9-]+)/uevent")
	match := re.FindStringSubmatch(path)
	return match[re.SubexpIndex("dev")]
}

func createDeviceFromUevent(path string) (*Device, error) {
	contents, err := readFile(path)
	if err != nil {
		return nil, err
	}

	res := new(Device)
	res.Name = getDeviceNameFromUevent(path)

	re := regexp.MustCompile("^?(?P<key>[A-Za-z]+)=(?P<value>[A-Za-z0-9_/]+)$?")
	matches := re.FindAllSubmatch(contents, -1)
	for _, match := range matches {
		key := string(match[re.SubexpIndex("key")])
		value := string(match[re.SubexpIndex("value")])

		switch strings.ToLower(key) {
		case "devnum":
			res.DevNum = value
		case "busnum":
			res.BusNum = value
		case "devname":
			res.DevName = value
		case "major":
			res.Major = value
		case "minor":
			res.Minor = value
		case "devpath":
			res.DevPath = filepath.Join("/dev", value)
		}
	}

	return res, nil
}

func readFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	contents := make([]byte, 256)
	f.Read(contents)
	return contents, nil
}

func getSerial(d *Device) string {
	f := fmt.Sprintf("/sys/bus/usb/devices/%s/serial", d.Name)
	contents, err := readFile(f)
	if err != nil {
		return ""
	}

	return string(contents)
}

func getProductId(d *Device) string {
	f := fmt.Sprintf("/sys/bus/usb/devices/%s/idProduct", d.Name)
	contents, err := readFile(f)
	if err != nil {
		return ""
	}

	return string(contents)
}

func getVendorId(d *Device) string {
	f := fmt.Sprintf("/sys/bus/usb/devices/%s/idVendor", d.Name)
	contents, err := readFile(f)
	if err != nil {
		return ""
	}

	return string(contents)
}

func getProduct(d *Device) string {
	f := fmt.Sprintf("/sys/bus/usb/devices/%s/product", d.Name)
	contents, err := readFile(f)
	if err != nil {
		return ""
	}

	return string(contents)
}

func getManufacturer(d *Device) string {
	f := fmt.Sprintf("/sys/bus/usb/devices/%s/manufacturer", d.Name)
	contents, err := readFile(f)
	if err != nil {
		return ""
	}

	return string(contents)
}
