package hid

import (
	"fmt"
	"sync"
	"time"
)

var JoystickReportDesc = []byte{
	0x05, 0x01, // USAGE_PAGE (Generic Desktop)
	0x09, 0x04, // USAGE (Joystick)
	0xa1, 0x01, // COLLECTION (Application)
	0x15, 0x81, //   LOGICAL_MINIMUM (-127)
	0x25, 0x7f, //   LOGICAL_MAXIMUM (127)
	0x05, 0x01, //   USAGE_PAGE (Generic Desktop)
	0x09, 0x01, //   USAGE (Pointer)

	0xa1, 0x00, //   COLLECTION (Physical)

	// USAGE AXIS 0-7
	0x09, 0x30,
	0x09, 0x31,
	0x09, 0x32,
	0x09, 0x33,
	0x09, 0x34,
	0x09, 0x35,
	0x09, 0x36,
	0x09, 0x37,

	0x75, 0x08, //     REPORT_SIZE (8)
	0x95, 0x08, //     REPORT_COUNT (8)
	0x81, 0x02, //     INPUT (Data,Var,Abs)
	0xc0, //   END_COLLECTION

	0x05, 0x09, //   USAGE_PAGE (Button)
	0x19, 0x01, //   USAGE_MINIMUM (Button 1)
	0x29, 0x0b, //   USAGE_MAXIMUM (Button 11)
	0x15, 0x00, //   LOGICAL_MINIMUM (0)
	0x25, 0x01, //   LOGICAL_MAXIMUM (1)
	0x75, 0x01, //   REPORT_SIZE (1)
	0x95, 0x0b, //   REPORT_COUNT (11)
	0x55, 0x00, //   UNIT_EXPONENT (0)
	0x65, 0x00, //   UNIT (None)
	0x81, 0x02, //   INPUT (Data,Var,Abs)

	// padding
	0x75, 0x01, //   REPORT_SIZE (1)
	0x95, 0x05, //   REPORT_COUNT (5)
	0x81, 0x02, //   INPUT (Data,Var,Abs)

	0xc0, // END_COLLECTION
}

type Joystick struct {
	*Device
	state     []byte
	sendMutex *sync.Mutex
}

func (j *Joystick) PressButton(num uint8) error {
	if num < 8 {
		j.state[8] |= byte(0x01) << num
	} else if num < 16 {
		j.state[9] |= byte(0x01) << (num - 8)
	}
	return j.SendState()
}

func (j *Joystick) ReleaseButton(num uint8) error {
	if num < 8 {
		j.state[8] &= ^(byte(0x01) << num)
	} else if num < 16 {
		j.state[9] &= ^(byte(0x01) << (num - 8))
	}
	return j.SendState()
}

func (j *Joystick) SetAxis(num uint8, value int16) error {
	if num >= 8 {
		return nil
	}
	prev := j.state[num]
	new := byte(value >> 8)
	if prev == new {
		return nil
	}
	j.state[num] = byte(value >> 8)
	return j.SendState()
}

func (j *Joystick) SendState() error {
	j.sendMutex.Lock()
	defer j.sendMutex.Unlock()
	err := j.Write(j.state)
	if err != nil {
		return fmt.Errorf("failed to set joystick state: %w", err)
	}
	return nil
}

func NewJoystick(path string) *Joystick {
	j := &Joystick{
		Device: &Device{
			path: path,
		},
		state:     make([]byte, 10),
		sendMutex: &sync.Mutex{},
	}

	// set initial trigger states
	j.state[2] = 0x80
	j.state[5] = 0x80

	// send initial state
	go func() {
		// wait for device to fully setup? somehow not working otherwise
		time.Sleep(time.Second)
		_ = j.SendState()
	}()

	return j
}
