package deck

import "os"

// SerialNumber requires root
func SerialNumber() (string, error) {
	serial, err := os.ReadFile("/sys/devices/virtual/nvme-subsystem/nvme-subsys0/serial")
	return string(serial), err
}
