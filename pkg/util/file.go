package util

import (
	"io"
	"os"
)

func CopyFile(in, out string) error {
	i, err := os.Open(in)
	if err != nil {
		return err
	}
	defer i.Close()
	o, err := os.Create(out)
	if err != nil {
		return err
	}
	defer o.Close()
	_, err = io.Copy(o, i)
	if err != nil {
		return err
	}
	err = o.Sync()
	if err != nil {
		return err
	}
	return nil
}
