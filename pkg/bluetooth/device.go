package bluetooth

import (
	"context"
	"fmt"
	"github.com/lucaber/deckjoy/pkg/hid"
	"github.com/lucaber/deckjoy/pkg/ipc"
)

type BluetoothDevice struct {
	Daemon ipc.DeckJoyDaemonClient
}

func NewBluetoothDevice(daemon ipc.DeckJoyDaemonClient) *BluetoothDevice {
	return &BluetoothDevice{daemon}
}

func (d *BluetoothDevice) Open() error {
	if d.Daemon == nil {
		return fmt.Errorf("no daemon connection")
	}
	return nil
}

func (d *BluetoothDevice) Write(data []byte) error {
	_, err := d.Daemon.WriteBluetoothHIDData(context.Background(), &ipc.WriteBluetoothHIDDataRequest{
		Data: data,
	})
	return err
}

var _ hid.Device = (*BluetoothDevice)(nil)
