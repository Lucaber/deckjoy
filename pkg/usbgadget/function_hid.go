package usbgadget

import (
	"errors"
	"os"
	"path"
	"strings"
)

type FunctionHID struct {
	*FunctionGeneric
}

type FunctionHIDAttributes struct {
	Protocol     uint8
	Subclass     uint8
	ReportLength uint8
	ReportDesc   []byte
}

func (g *Gadget) CreateFunctionHID(name string) (*FunctionHID, error) {
	gf := g.GetFunction(name)
	f := &FunctionHID{
		gf,
	}

	if err := mkDir(f.path); err != nil {
		return nil, err
	}

	return f, nil
}

func (f *FunctionHID) SetAttributes(attrs FunctionHIDAttributes) error {
	if err := writeHex8(path.Join(f.path, "protocol"), attrs.Protocol); err != nil {
		return err
	}
	if err := writeHex8(path.Join(f.path, "subclass"), attrs.Protocol); err != nil {
		return err
	}
	if err := writeHex8(path.Join(f.path, "report_length"), attrs.ReportLength); err != nil {
		return err
	}
	if err := writeBytes(path.Join(f.path, "report_desc"), attrs.ReportDesc); err != nil {
		return err
	}
	return nil
}

func (f *FunctionHID) GetDevicePath() (string, error) {
	devContent, err := os.ReadFile(path.Join(f.path, "dev"))
	if err != nil {
		return "", err
	}
	deviceNumbers := strings.Split(string(devContent), "\n")[0]

	uevent, err := os.ReadFile(path.Join("/sys/dev/char", deviceNumbers, "uevent"))
	if err != nil {
		return "", err
	}

	rows := strings.Split(string(uevent), "\n")
	for _, row := range rows {
		entry := strings.SplitN(row, "=", 2)
		if entry[0] == "DEVNAME" {
			return path.Join("/dev", entry[1]), nil
		}
	}

	return "", errors.New("could not find device path")
}

func (f *FunctionHID) SetDevicePermissions() error {
	p, err := f.GetDevicePath()
	if err != nil {
		return err
	}
	err = os.Chmod(p, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
