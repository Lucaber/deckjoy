package hid

func AddReportID(descriptor []byte, id byte) []byte {
	out := make([]byte, 0, len(descriptor)+2)
	out = append(out, descriptor[:6]...)
	out = append(out, 0x85, id)
	out = append(out, descriptor[6:]...)
	return out
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
