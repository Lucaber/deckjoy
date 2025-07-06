package bluetooth

import (
	"context"
	"errors"
	"github.com/godbus/dbus/v5"
	"github.com/muka/go-bluetooth/bluez/profile/profile"
	log "github.com/sirupsen/logrus"
)

type Bluetooth struct {
	sockets []*Socket
}

func NewBluetooth() *Bluetooth {
	return &Bluetooth{}
}

func (b *Bluetooth) Run(ctx context.Context, record SDPRecord) error {
	profileManager1, err := profile.NewProfileManager1()
	if err != nil {
		return err
	}

	path := dbus.ObjectPath("/bluez/lucaber/deckjoy")
	profile1, err := profile.NewProfile1(string(path), path)
	if err != nil {
		return err
	}
	log.Debugf("bluetooth profile1: %+v\n", profile1)
	ch, _, err := profile1.GetObjectManagerSignal()
	if err != nil {
		return err
	}

	uuid := "00001124-0000-1000-8000-00805f9b34fb"
	err = profileManager1.RegisterProfile(path, uuid, map[string]interface{}{
		"ServiceRecord": string(record.String()),
		//"RequireAuthentication": false,
		//"RequireAuthorization":  false,
		"AutoConnect": true,
		//"Role":                  "server",
	})
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case ev := <-ch:
				log.Infof("bluetooth event: %+v\n", ev)
			}
		}
	}()

	c, err := NewL2CAPListener(0x11)
	if err != nil {
		return err
	}
	s, err := NewL2CAPListener(0x13)
	if err != nil {
		return err
	}

	go func() {
		for {
			_, err := c.Accept()
			if err != nil {
				log.Errorf("bluetooth accept err 0x11: %+v\n", err)
			}
			d, err := s.Accept()
			if err != nil {
				log.Errorf("bluetooth accept err 0x13: %+v\n", err)
			}
			b.sockets = append(b.sockets, d)
			log.Infof("new device connected")
		}
	}()

	return nil
}

func (b *Bluetooth) Write(data []byte) error {
	if len(b.sockets) == 0 {
		return errors.New("not connected")
	}
	ok := false
	var lasterr error
	for _, d := range b.sockets {
		_, err := d.Write(data)
		if err != nil {
			lasterr = err
		} else {
			ok = true
		}
	}

	if !ok {
		return lasterr
	}
	return nil
}
