package bluetooth

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
	"sync"
)

type Listener struct {
	fd int
}

func NewL2CAPListener(psm uint16) (*Listener, error) {
	var err error
	fd, err := unix.Socket(unix.AF_BLUETOOTH, unix.SOCK_SEQPACKET, unix.BTPROTO_L2CAP)
	if err != nil {
		return nil, errors.Wrap(err, "can't create socket")
	}

	if err := unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_REUSEADDR, 1); err != nil {
		return nil, errors.Wrap(err, "can't set socket flags")
	}

	if err := unix.Bind(fd, &unix.SockaddrL2{
		PSM:      psm,
		Addr:     [6]uint8{},
		CID:      0,
		AddrType: 0,
	}); err != nil {
		return nil, errors.Wrap(err, "can't bind")
	}

	if err := unix.Listen(fd, 1); err != nil {
		return nil, errors.Wrap(err, "can't listen")
	}

	return &Listener{fd: fd}, nil
}

func (l *Listener) Accept() (*Socket, error) {
	fd, _, err := unix.Accept(l.fd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to accept")
	}
	return &Socket{fd: fd, lock: &sync.Mutex{}}, nil
}

type Socket struct {
	fd   int
	lock *sync.Mutex
}

func (s *Socket) Write(p []byte) (int, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	n, err := unix.Write(s.fd, p)
	log.Debugf("wrote %v", p)
	return n, errors.Wrap(err, "can't write socket")
}

func (s *Socket) Close() error {
	return errors.Wrap(unix.Close(s.fd), "can't close socket")
}
