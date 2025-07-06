package main

import (
	"github.com/lucaber/deckjoy/pkg/bluetooth"
	"github.com/lucaber/deckjoy/pkg/cmd"
	"github.com/lucaber/deckjoy/pkg/hid"
	"github.com/lucaber/deckjoy/pkg/usb"
	"github.com/urfave/cli/v2"
	log "github.com/sirupsen/logrus"
	"os"
	"runtime"
)

func main() {
	runtime.LockOSThread()
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:   "daemon",
				Usage:  "daemon to configure the steam deck; start with root permissions",
				Action: cmd.RunDaemon,
			},
			{
				Name:  "cleanup",
				Usage: "remove usb gadget",
				Action: func(*cli.Context) error {
					deckUSB, err := usb.NewUSB()
					if err != nil {
						return err
					}
					err = deckUSB.Destroy()
					if err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:   "gui",
				Usage:  "show gui",
				Action: cmd.RunGui,
			},
			{
				Name:  "test",
				Usage: "test",
				Action: func(context *cli.Context) error {
					log.SetLevel(log.DebugLevel)
					err := bluetooth.NewBluetooth().Run(context.Context, bluetooth.SDPRecord{HIDDescriptor: hid.KeyboardReportDesc})
					if err != nil {
						return err
					}
					select {}
				},
			},
		},
		Action: cmd.RunGui,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
