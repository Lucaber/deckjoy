package gui

import (
	"context"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/lucaber/deckjoy/pkg/config"
	"github.com/lucaber/deckjoy/pkg/service"
	"github.com/lucaber/deckjoy/pkg/usbgadget"
	log "github.com/sirupsen/logrus"
	"golang.org/x/image/colornames"
	"time"
)

type GUI struct {
	deck        *service.Deck
	app         fyne.App
	window      fyne.Window
	inputWindow *InputWindow
}

func NewGUI(deck *service.Deck) *GUI {
	return &GUI{
		deck: deck,
	}
}

func (g *GUI) Run() {
	g.app = app.New()
	g.window = g.app.NewWindow("DeckJoy")
	g.window.SetFullScreen(true)

	errLabel := widget.NewLabel("")
	startUSBButton := widget.NewButton("Start USB", func() {
		g.deck.StartDaemon(context.Background())
	})
	startInputWindowButton := widget.NewButton("Show Mouse/Keyboard", func() {
		go func() {
			if g.inputWindow == nil {
				g.inputWindow = NewInputWindow(g.deck)
				err := g.inputWindow.Run()
				if err != nil {
					log.WithError(err).Errorf("input window crashed")
				}
			} else {
				g.inputWindow.Show()
			}
		}()
	})
	startInputWindowButton.Disable()

	// todo: sometimes not working, maybe before libcomposite is loaded?
	if hasDevice, err := usbgadget.HasDevice(); err != nil || !hasDevice {
		startUSBButton.Disable()
		startUSBButton.SetText("No USB port in DRD mode found")
	}

	go func() {
		// wait for daemon to start
		for {
			time.Sleep(100 * time.Millisecond)
			if g.deck.Mouse != nil {
				break
			}
			if g.deck.SetupErr != nil {
				errLabel.SetText(g.deck.SetupErr.Error())
			}
		}
		startUSBButton.Disable()
		startInputWindowButton.Enable()
	}()

	versionText := canvas.NewText(fmt.Sprintf("%s", config.Version), colornames.Gray)
	versionText.TextSize = 8

	g.window.SetContent(container.NewVBox(
		widget.NewLabel(fmt.Sprintf("DeckJoy")),
		startUSBButton,
		startInputWindowButton,
		errLabel,
		layout.NewSpacer(),
		container.NewHBox(
			layout.NewSpacer(),
			versionText,
		),
	))

	g.window.ShowAndRun()
}
