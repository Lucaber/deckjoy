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

	g.setBlackScreen(false)
	g.window.ShowAndRun()
}

func (g *GUI) setBlackScreen(on bool) {
	g.window.SetPadded(!on)
	if on {
		g.window.SetContent(g.buildBlackScreenUI())
		return
	}
	g.window.SetContent(g.buildMainUI())
}

func (g *GUI) buildMainUI() fyne.CanvasObject {
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

	blackScreenButton := newHiddenCursorButton("Black Screen", func() {
		g.setBlackScreen(true)
	})

	return container.NewVBox(
		widget.NewLabel(fmt.Sprintf("DeckJoy")),
		startUSBButton,
		startInputWindowButton,
		blackScreenButton,
		errLabel,
		layout.NewSpacer(),
		container.NewHBox(
			layout.NewSpacer(),
			versionText,
		),
	)
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

func (g *GUI) buildBlackScreenUI() fyne.CanvasObject {
	return newTappableBlackRect(func() {
		g.setBlackScreen(false)
	})
}

// tappableBlackRect is a full-screen black widget that restores the main UI on any tap/touch.
// It hides the mouse cursor immediately and fills the window with pure black.
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
	rect := canvas.NewRectangle(colornames.Black)
	return &blackRectRenderer{rect: rect}
}

func (r *tappableBlackRect) Tapped(_ *fyne.PointEvent) {
	if r.onTap != nil {
		r.onTap()
	}
}

func (r *tappableBlackRect) TappedSecondary(_ *fyne.PointEvent) {}

// Cursor implements desktop.Cursorable — hides the mouse pointer immediately.
func (r *tappableBlackRect) Cursor() desktop.Cursor { return desktop.HiddenCursor }

func (r *tappableBlackRect) MouseIn(_ *desktop.MouseEvent)    {}
func (r *tappableBlackRect) MouseMoved(_ *desktop.MouseEvent) {}
func (r *tappableBlackRect) MouseOut()                        {}

type blackRectRenderer struct {
	rect *canvas.Rectangle
}

func (r *blackRectRenderer) Layout(size fyne.Size) {
	r.rect.Move(fyne.NewPos(0, 0))
	r.rect.Resize(size)
}
func (r *blackRectRenderer) MinSize() fyne.Size           { return fyne.NewSize(0, 0) }
func (r *blackRectRenderer) Refresh()                     { r.rect.Refresh() }
func (r *blackRectRenderer) BackgroundColor() color.Color { return color.Black }
func (r *blackRectRenderer) Objects() []fyne.CanvasObject { return []fyne.CanvasObject{r.rect} }
func (r *blackRectRenderer) Destroy()                     {}
