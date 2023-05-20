package hid

import (
	"errors"
	"os"
)

type Device struct {
	path string
	file *os.File
}

func (d *Device) Open() error {
	file, err := os.OpenFile(d.path, os.O_RDWR, os.ModeCharDevice)
	if err != nil {
		return err
	}
	d.file = file
	return nil
}

func (d *Device) Write(b []byte) error {
	if d.file == nil {
		if err := d.Open(); err != nil {
			return err
		}
	}

	if _, err := d.file.Write(b); err != nil && !errors.Is(err, os.ErrClosed) {
		return err
	} else if err != nil {
		return d.Write(b)
	}
	return nil
}
