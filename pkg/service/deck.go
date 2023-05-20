package service

import (
	"context"
	"fmt"
	"github.com/lucaber/deckjoy/pkg/daemon"
	"github.com/lucaber/deckjoy/pkg/hid"
	"github.com/lucaber/deckjoy/pkg/ipc"
	"github.com/lucaber/deckjoy/pkg/joystick"
	log "github.com/sirupsen/logrus"
	"time"
)

const SocketPath = "/run/deckjoy.sock"

type Deck struct {
	Keyboard      *hid.Keyboard
	Joystick      *hid.Joystick
	Mouse         *hid.Mouse
	LocalJoystick *joystick.Joystick
	Daemon        ipc.DeckJoyDaemonClient
	daemonStop    context.CancelFunc
	SetupErr      error
}

func NewDeck() *Deck {
	deck := &Deck{
		LocalJoystick: joystick.NewJoystick("/dev/input/js0"),
	}
	return deck
}

func (d *Deck) Run(ctx context.Context) {
	for {
		var err error
		connectCtx, connectCtxCancel := context.WithTimeout(ctx, time.Second)
		d.Daemon, err = daemon.NewClient(connectCtx, SocketPath)
		connectCtxCancel()
		if err != nil {
			log.WithError(err).Infof("failed to connect to daemon")
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	_, err := d.Daemon.Init(ctx, &ipc.Empty{})
	if err != nil {
		d.SetupErr = fmt.Errorf("usb init failed: %w", err)
		log.WithError(err).Infof("usb init failed")
		return
	}

	joystickRes, err := d.Daemon.SetupJoystick(ctx, &ipc.SetupJoystickRequest{
		UserPermissions: true,
	})
	if err != nil {
		d.SetupErr = fmt.Errorf("joystick setup failed: %w", err)
		log.WithError(err).Infof("joystick setup failed")
		return
	}
	log.Infof("created joystick at %s", joystickRes.Path)
	d.Joystick = hid.NewJoystick(joystickRes.Path)

	keyboardRes, err := d.Daemon.SetupKeyboard(ctx, &ipc.SetupKeyboardRequest{
		UserPermissions: true,
	})
	if err != nil {
		d.SetupErr = fmt.Errorf("keyboard setup failed: %w", err)
		log.WithError(err).Infof("keyboard setup failed")
		return
	}
	log.Infof("created keyboard at %s", keyboardRes.Path)
	d.Keyboard = hid.NewKeyboard(keyboardRes.Path)

	mouseRes, err := d.Daemon.SetupMouse(ctx, &ipc.SetupMouseRequest{
		UserPermissions: true,
	})
	if err != nil {
		d.SetupErr = fmt.Errorf("mouse setup failed: %w", err)
		log.WithError(err).Infof("mouse setup failed")
		return
	}
	log.Infof("created mouse at %s", mouseRes.Path)
	d.Mouse = hid.NewMouse(mouseRes.Path)

	d.RunJoystick()
}

func (d *Deck) StartDaemon(ctx context.Context) {
	daemonCtx, cancel := context.WithCancel(ctx)
	d.daemonStop = cancel
	errors := daemon.RunDaemonProcess(daemonCtx)
	go func() {
		for e := range errors {
			log.WithError(e).Error("daemon run error")
		}
		d.daemonStop = nil
		log.Infof("daemon exited")
	}()
}

func (d *Deck) Stop() {
	d.StopDaemon()
}

func (d *Deck) StopDaemon() {
	if d.daemonStop == nil {
		// only stop daemon if we started it
		return
	}

	stopCtx, stopCtxCancel := context.WithTimeout(context.Background(), time.Second)
	log.Infof("stopping daemon")
	err := daemon.StopDaemonProcess(stopCtx, SocketPath)
	if err != nil {
		log.WithError(err).Errorf("failed to stop daemon")
	}
	// wait for process to exit normally
	time.Sleep(time.Second)

	stopCtxCancel()
	if d.daemonStop != nil {
		d.daemonStop()
	}
}
