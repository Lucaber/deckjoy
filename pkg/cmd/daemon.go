package cmd

import (
	"github.com/lucaber/deckjoy/pkg/daemon"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func RunDaemon(*cli.Context) error {
	d := daemon.NewServer("/run/deckjoy.sock")

	stopOnce := &sync.Once{}
	stopFunc := func() {
		log.Infof("closing server")
		err := d.Close()
		if err != nil {
			log.WithError(err).Error("failed to close server")
		}
		log.Infof("closed server")
	}
	defer stopOnce.Do(stopFunc)

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		stopOnce.Do(stopFunc)
		os.Exit(0)
	}()

	err := d.Run()
	if err != nil {
		return err
	}
	return nil
}
