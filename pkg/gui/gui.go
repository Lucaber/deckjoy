package gui

import (
	"context"
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/lucaber/deckjoy/pkg/config"
	"github.com/lucaber/deckjoy/pkg/service"
	log "github.com/sirupsen/logrus"
	"golang.org/x/image/colornames"
)

type GUI struct {
	deck        *service.Deck
	app         fyne.App
	window      fyne.Window
	inputWindow *InputWindow
	mainContent fyne.CanvasObject
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
	errLabel.Wrapping = fyne.TextWrapWord
	startInputWindowButton := widget.NewButton("Show Mouse/Keyboard", func() {
		g.showInputWindow(false)
	})
	startBlackScreenButton := newHiddenCursorButton("Black Screen", func() {
		g.setFyneBlackScreen()
	})
	startInputWindowButton.Disable()
	startBlackScreenButton.Disable()

	var startUSBButton *widget.Button
	startUSBButton = widget.NewButton("Start USB", func() {
		startUSBButton.Disable()
		errLabel.SetText("")

		g.deck.StartDaemon(context.Background())

		go func() {
			for {
				time.Sleep(100 * time.Millisecond)
				if g.deck.Mouse != nil {
					startInputWindowButton.Enable()
					startBlackScreenButton.Enable()
					return
				}
				if g.deck.SetupErr != nil {
					errLabel.SetText(g.deck.SetupErr.Error())
					startUSBButton.Enable()
					return
				}
			}
		}()
	})

	versionText := canvas.NewText(fmt.Sprintf("%s", config.Version), colornames.Gray)
	versionText.TextSize = 8

	g.mainContent = container.NewVBox(
		widget.NewLabel(fmt.Sprintf("DeckJoy")),
		startUSBButton,
		startInputWindowButton,
		startBlackScreenButton,
		errLabel,
		layout.NewSpacer(),
		container.NewHBox(
			layout.NewSpacer(),
			versionText,
		),
	)
	g.window.SetContent(g.mainContent)

	g.window.ShowAndRun()
}

func (g *GUI) setFyneBlackScreen() {
	g.window.SetPadded(false)
	blackRect := newTappableBlackRect(func() {
		g.window.SetPadded(true)
		g.window.SetContent(g.mainContent)
	})
	g.window.SetContent(blackRect)
}

type hiddenCursorButton struct {
	widget.Button
}

func newHiddenCursorButton(label string, tapped func()) *hiddenCursorButton {
	b := &hiddenCursorButton{}
	b.Text = label
	b.OnTapped = tapped
	b.ExtendBaseWidget(b)
	return b
}

func (b *hiddenCursorButton) Cursor() desktop.Cursor { return desktop.HiddenCursor }

// tappableBlackRect is a full-screen black widget that restores the main UI on tap.
type tappableBlackRect struct {
	widget.BaseWidget
	onTap func()
}

func newTappableBlackRect(onTap func()) *tappableBlackRect {
	r := &tappableBlackRect{onTap: onTap}
	r.ExtendBaseWidget(r)
	return r
}

func (r *tappableBlackRect) CreateRenderer() fyne.WidgetRenderer {
	rect := canvas.NewRectangle(color.Black)
	return widget.NewSimpleRenderer(rect)
}

func (r *tappableBlackRect) Tapped(_ *fyne.PointEvent) {
	if r.onTap != nil {
		r.onTap()
	}
}

func (r *tappableBlackRect) TappedSecondary(_ *fyne.PointEvent) {}

func (r *tappableBlackRect) Cursor() desktop.Cursor { return desktop.HiddenCursor }

func (r *tappableBlackRect) MouseIn(_ *desktop.MouseEvent)    {}
func (r *tappableBlackRect) MouseMoved(_ *desktop.MouseEvent) {}
func (r *tappableBlackRect) MouseOut()                        {}

func (g *GUI) showInputWindow(blackScreen bool) {
	go func() {
		if g.inputWindow == nil {
			g.inputWindow = NewInputWindow(g.deck)
			g.inputWindow.blackScreen = blackScreen
			err := g.inputWindow.Run()
			if err != nil {
				log.WithError(err).Errorf("input window crashed")
			}
		} else {
			g.inputWindow.blackScreen = blackScreen
			g.inputWindow.Show()
		}
	}()
}
