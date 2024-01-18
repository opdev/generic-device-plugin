package rm

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func findDevices() ([]Device, error) {
	uevents, err := filepath.Glob("/sys/bus/usb/devices/[0-9]-[0-9]/uevent")
	if err != nil {
		return nil, err
	}

	var devices []Device
	for _, device := range uevents {
		dev, _ := parseDeviceInfo(device)
		devices = append(devices, *dev)
	}

	return devices, nil
}

func parseDeviceInfo(uevent string) (*Device, error) {
	device, err := createDeviceFromUevent(uevent)
	if err != nil {
		return nil, err
	}

	// Set properties not in the uevent file
	setSerial(device)
	setProductId(device)
	setVendorId(device)

	return device, nil
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

func setSerial(d *Device) {
	f := fmt.Sprintf("/sys/bus/usb/devices/%s/serial", d.Name)
	contents, err := readFile(f)
	if err != nil {
		return
	}

	d.Serial = string(contents)
}

func setProductId(d *Device) {
	f := fmt.Sprintf("/sys/bus/usb/devices/%s/idProduct", d.Name)
	contents, err := readFile(f)
	if err != nil {
		return
	}

	d.ProductId = string(contents)
}

func setVendorId(d *Device) {
	f := fmt.Sprintf("/sys/bus/usb/devices/%s/idVendor", d.Name)
	contents, err := readFile(f)
	if err != nil {
		return
	}

	d.VenderId = string(contents)
}
