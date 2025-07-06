package hid

import (
	"errors"
	"os"
)

type Device interface {
	Open() error
	Write(b []byte) error
}

type FileDevice struct {
	Path string
	file *os.File
}

func NewFileDevice(path string) *FileDevice {
	return &FileDevice{path, nil}
}

func (d *FileDevice) Open() error {
	file, err := os.OpenFile(d.Path, os.O_RDWR, os.ModeCharDevice)
	if err != nil {
		return err
	}
	d.file = file
	return nil
}

func (d *FileDevice) Write(b []byte) error {
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

type NullDevice struct {
}

func (n NullDevice) Open() error {
	return nil
}

func (n NullDevice) Write(b []byte) error {
	return nil
}

func NewNullDevice() *NullDevice {
	return &NullDevice{}
}

var _ Device = &NullDevice{}
