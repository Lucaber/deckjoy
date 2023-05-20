package usbgadget

import (
	"errors"
	"os"
	"strconv"
)

func mkDir(path string) error {
	return os.MkdirAll(path, os.ModeDir)
}

func writeString(path string, s string) error {
	return os.WriteFile(path, []byte(s), os.ModePerm)
}

func writeHex16(path string, i uint16) error {
	hexStr := strconv.FormatInt(int64(i), 16)
	if len(hexStr) > 4 {
		return errors.New("WTF")
	}
	for len(hexStr) < 4 {
		hexStr = "0" + hexStr
	}
	data := "0x" + hexStr
	return os.WriteFile(path, []byte(data), os.ModePerm)
}

func writeHex8(path string, i uint8) error {
	hexStr := strconv.FormatInt(int64(i), 16)
	if len(hexStr) > 2 {
		return errors.New("WTF")
	}
	for len(hexStr) < 2 {
		hexStr = "0" + hexStr
	}
	data := "0x" + hexStr
	return os.WriteFile(path, []byte(data), os.ModePerm)
}

func writeBytes(path string, i []byte) error {
	return os.WriteFile(path, i, os.ModePerm)
}

func listDir(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(entries))
	for i, e := range entries {
		names[i] = e.Name()
	}
	return names, nil
}
