package hid

func AddReportID(descriptor []byte, id byte) []byte {
	x := descriptor[:6]
	x = append(x, 0x85)
	x = append(x, id)
	x = append(x, descriptor[6:]...)
	return x
}

type ReportIDDevice struct {
	inner Device
	id    byte
}

func (r *ReportIDDevice) Open() error {
	return r.inner.Open()
}

func (r *ReportIDDevice) Write(b []byte) error {
	return r.inner.Write(append([]byte{0xA1, r.id}, b...))
}

var _ Device = &ReportIDDevice{}

func NewReportIDDevice(inner Device, id byte) *ReportIDDevice {
	return &ReportIDDevice{
		inner: inner,
		id:    id,
	}
}
