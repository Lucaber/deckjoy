package usbgadget

func ListDevices() ([]string, error) {
	return listDir("/sys/class/udc")
}

func HasDevice() (bool, error) {
	entries, err := ListDevices()
	if err != nil {
		return false, err
	}

	return len(entries) > 0, nil
}
