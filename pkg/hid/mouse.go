package hid

var MouseReportDesc = []byte{
	0x05, 0x01, // USAGE_PAGE (Generic Desktop)
	0x09, 0x02, // USAGE (Mouse)
	0xa1, 0x01, // COLLECTION (Application)
	0x09, 0x01, // USAGE (Pointer)
	0xa1, 0x00, // COLLECTION (Physical)
	0x05, 0x09, // USAGE_PAGE (Button)
	0x19, 0x01, // USAGE_MINIMUM (Button 1)
	0x29, 0x03, // USAGE_MAXIMUM (Button 3)
	0x15, 0x00, // LOGICAL_MINIMUM (0)
	0x25, 0x01, // LOGICAL_MAXIMUM (1)
	0x95, 0x03, // REPORT_COUNT (3)
	0x75, 0x01, // REPORT_SIZE (1)
	0x81, 0x02, // INPUT (Data,Var,Abs)
	0x95, 0x01, // REPORT_COUNT (1)
	0x75, 0x05, // REPORT_SIZE (5)
	0x81, 0x03, // INPUT (Cnst,Var,Abs)
	0x05, 0x01, // USAGE_PAGE (Generic Desktop)
	0x09, 0x30, // USAGE (X)
	0x09, 0x31, // USAGE (Y)
	0x16, 0x01, 0x80, //     Logical Minimum (-32767)
	0x26, 0xFF, 0x7F, //     Logical Maximum (32767)
	//0x15, 0x81, // LOGICAL_MINIMUM (-127)
	//0x25, 0x7f, // LOGICAL_MAXIMUM (127)
	0x75, 0x10, // REPORT_SIZE (16)
	0x95, 0x02, // REPORT_COUNT (2)
	0x81, 0x06, // INPUT (Data,Var,Rel)
	0xc0, // END_COLLECTION
	0xc0, // END_COLLECTION
}

type Mouse struct {
	*Device
	buttons byte
}

func (k *Mouse) Move(x, y int16) error {
	err := k.Write([]byte{k.buttons,
		byte(x & 0xff),
		byte(x >> 8),
		byte(y & 0xff),
		byte(y >> 8),
	})
	if err != nil {
		return err
	}
	return k.Write([]byte{k.buttons, 0x00, 0x00, 0x00, 0x00})
}

type MouseButton uint8

const MouseButtonLeft = MouseButton(0)
const MouseButtonRight = MouseButton(1)
const MouseButtonMiddle = MouseButton(2)

func (k *Mouse) PressButton(button MouseButton) error {
	k.buttons |= 0x1 << button
	return k.Write([]byte{k.buttons, 0x00, 0x00, 0x00, 0x00})
}
func (k *Mouse) ReleaseButton(button MouseButton) error {
	k.buttons &= ^(0x1 << button)
	return k.Write([]byte{k.buttons, 0x00, 0x00, 0x00, 0x00})
}

func NewMouse(path string) *Mouse {
	k := &Mouse{
		Device: &Device{
			path: path,
		},
		buttons: 0x00,
	}
	return k
}
