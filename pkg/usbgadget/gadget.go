package usbgadget

import (
	"errors"
	"os"
	"path"
)

type Gadget struct {
	path string
}

type GadgetAttributes struct {
	IDVendor  uint16
	IDProduct uint16
	BCDDevice uint16
	BCDUSB    uint16
}
type GadgetStrings struct {
	SerialNumber string
	Manufacturer string
	Product      string
}

func CreateGadget(configfs string, name string) (*Gadget, error) {
	g := &Gadget{
		path: path.Join(configfs, "usb_gadget", name),
	}

	if err := mkDir(g.path); err != nil {
		return nil, err
	}

	return g, nil
}

func (g *Gadget) SetAttributes(attrs GadgetAttributes) error {
	if err := writeHex16(path.Join(g.path, "idVendor"), attrs.IDVendor); err != nil {
		return err
	}
	if err := writeHex16(path.Join(g.path, "idProduct"), attrs.IDProduct); err != nil {
		return err
	}
	if err := writeHex16(path.Join(g.path, "bcdDevice"), attrs.BCDDevice); err != nil {
		return err
	}
	if err := writeHex16(path.Join(g.path, "bcdUSB"), attrs.BCDUSB); err != nil {
		return err
	}

	return nil
}

func (g *Gadget) SetStrings(strs GadgetStrings) error {
	stringsPath := path.Join(g.path, "strings", "0x409")
	_ = mkDir(stringsPath)

	if err := writeString(path.Join(stringsPath, "serialnumber"), strs.SerialNumber); err != nil {
		return err
	}
	if err := writeString(path.Join(stringsPath, "manufacturer"), strs.Manufacturer); err != nil {
		return err
	}
	if err := writeString(path.Join(stringsPath, "product"), strs.Product); err != nil {
		return err
	}

	return nil
}

func (g *Gadget) DeleteStrings() error {
	stringsPath := path.Join(g.path, "strings", "0x409")
	return os.Remove(stringsPath)
}

func (g *Gadget) Enable() error {
	entries, err := ListDevices()
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		return errors.New("no udc device found")
	}

	dev := entries[0]
	return g.EnableUDC(dev)
}

func (g *Gadget) EnableUDC(udc string) error {
	return writeString(path.Join(g.path, "UDC"), udc)
}

func (g *Gadget) Disable() error {
	// writing "" isnt as reliable
	//return writeString(path.Join(g.path, "UDC"), "")
	return writeString(path.Join(g.path, "UDC"), "\n")
}

func (g *Gadget) Destroy(skipErrors bool) error {
	if err := g.Disable(); err != nil {
		if !skipErrors {
			return err
		}
	}

	configs, err := g.GetConfigs()
	if err == nil {
		for _, config := range configs {
			err = config.Destroy(skipErrors)
			if err != nil {
				return err
			}
		}
	} else {
		if !skipErrors {
			return err
		}
	}

	functions, err := g.GetFunctions()
	if err == nil {
		for _, function := range functions {
			err = function.Destroy(skipErrors)
			if err != nil {
				return err
			}
		}
	} else {
		if !skipErrors {
			return err
		}
	}

	if err := g.DeleteStrings(); err != nil {
		if !skipErrors {
			return err
		}
	}
	return os.Remove(g.path)
}
