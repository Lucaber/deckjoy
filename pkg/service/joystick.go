package service

import (
	"github.com/lucaber/deckjoy/pkg/joystick"
	log "github.com/sirupsen/logrus"
)

func (d *Deck) RunJoystick() {
	go func() {
		for {
			err := d.LocalJoystick.Open()
			if err != nil {
				log.WithError(err).Errorf("failed to open joystick")
				continue
			}

			events := d.LocalJoystick.Run()
			for ev := range events {
				switch ev.Type {
				case joystick.JS_EVENT_BUTTON:
					log.WithField("button", ev.Number).Debugf("received button event")
					if ev.Value != 0 {
						err = d.Joystick.PressButton(ev.Number)
						if err != nil {
							log.WithError(err).WithField("button", ev.Number).Errorf("failed to press button")
						}
					} else {
						err = d.Joystick.ReleaseButton(ev.Number)
						if err != nil {
							log.WithError(err).WithField("button", ev.Number).Errorf("failed to release button")
						}
					}
				case joystick.JS_EVENT_AXIS:
					log.WithField("axis", ev.Number).WithField("value", ev.Value).Debugf("received update axis event")
					err = d.Joystick.SetAxis(ev.Number, ev.Value)
					if err != nil {
						log.WithError(err).WithField("axis", ev.Number).WithField("value", ev.Value).Errorf("failed to update axis")
					}
				}
			}
		}
	}()
}
