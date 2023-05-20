package joystick

// linux/joystick.h

const JS_EVENT_BUTTON = 0x01 /* button pressed/released */
const JS_EVENT_AXIS = 0x02   /* joystick moved */
const JS_EVENT_INIT = 0x80   /* initial state of device */

type js_event struct {
	Time   uint32 /* event timestamp in milliseconds */
	Value  int16  /* value */
	Type   uint8  /* event type */
	Number uint8  /* axis/button number */
}

var JSIOCGAXES = _IOR('j', 0x11, 1)    /* get number of axes */
var JSIOCGBUTTONS = _IOR('j', 0x12, 1) /* get number of buttons */

// asm/ioctl.h
const (
	_IOC_NRBITS   = 8
	_IOC_TYPEBITS = 8
	_IOC_SIZEBITS = 14
	_IOC_DIRBITS  = 2

	_IOC_NRSHIFT   = 0
	_IOC_TYPESHIFT = (_IOC_NRSHIFT + _IOC_NRBITS)
	_IOC_SIZESHIFT = (_IOC_TYPESHIFT + _IOC_TYPEBITS)
	_IOC_DIRSHIFT  = (_IOC_SIZESHIFT + _IOC_SIZEBITS)

	_IOC_WRITE = 1
	_IOC_READ  = 2
)

func _IOC(dir int, t int, nr int, size int) uint {
	return uint((dir << _IOC_DIRSHIFT) | (t << _IOC_TYPESHIFT) | (nr << _IOC_NRSHIFT) | (size << _IOC_SIZESHIFT))
}

func _IOR(t int, nr int, size int) uint {
	return _IOC(_IOC_READ, t, nr, size)
}

func _IOW(t int, nr int, size int) uint {
	return _IOC(_IOC_WRITE, t, nr, size)
}
