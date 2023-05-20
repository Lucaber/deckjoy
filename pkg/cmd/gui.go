package cmd

import (
	"github.com/lucaber/deckjoy/pkg/gui"
	"github.com/lucaber/deckjoy/pkg/service"
	"github.com/lucaber/deckjoy/pkg/setup"
	"github.com/lucaber/deckjoy/pkg/steamworks"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func RunGui(c *cli.Context) error {

	err := setup.Install()
	if err != nil {
		log.WithError(err).Error("installation failed")
	}

	// requires a "steam_appid.txt" file in the root directory
	// but then steam drm blocks the steam account from starting a different game on a different device
	// with error: "another computer already playing"
	err = steamworks.Init()
	if err != nil {
		log.WithError(err).Error("failed to init Steamworks")
	}

	deck := service.NewDeck()
	go deck.Run(c.Context)

	stopOnce := &sync.Once{}
	stopFunc := func() {
		deck.Stop()
	}
	defer stopOnce.Do(stopFunc)

	signals := make(chan os.Signal, 2)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signals
		stopOnce.Do(stopFunc)
		os.Exit(0)
	}()

	gui := gui.NewGUI(deck)
	gui.Run()
	return nil
}
