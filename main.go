package main

import (
	"github.com/lucaber/deckjoy/pkg/cmd"
	"github.com/lucaber/deckjoy/pkg/setup"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
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
					deckSetup, err := setup.NewDeck()
					if err != nil {
						return err
					}
					err = deckSetup.Destroy()
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
		},
		Action: cmd.RunGui,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
