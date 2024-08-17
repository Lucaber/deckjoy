package gui

import (
	"errors"
	"fmt"
	"github.com/lucaber/deckjoy/pkg/hid"
	"github.com/lucaber/deckjoy/pkg/service"
	"github.com/lucaber/deckjoy/pkg/steamworks"
	log "github.com/sirupsen/logrus"
	"runtime"
	"sync"
)
import "github.com/veandco/go-sdl2/sdl"

type InputWindow struct {
	deck                    *service.Deck
	window                  *sdl.Window
	renderer                *sdl.Renderer
	keyboardGui             *KeyboardGUI
	pendingEventsBatch      []sdl.Event
	pendingEventsBatchMutex *sync.Mutex
	// Unlock to trigger event handling
	pendingEventsLoopLock *sync.Mutex
	touches               map[sdl.FingerID]*Key
}

var QuitErr = fmt.Errorf("quit")

func NewInputWindow(deck *service.Deck) *InputWindow {
	return &InputWindow{
		deck:                    deck,
		pendingEventsBatch:      []sdl.Event{},
		touches:                 map[sdl.FingerID]*Key{},
		pendingEventsBatchMutex: &sync.Mutex{},
		pendingEventsLoopLock:   &sync.Mutex{},
	}
}

func (iw *InputWindow) Run() error {
	runtime.LockOSThread()
	go iw.runPendingEventsHandler()

	var err error
	sdl.Main(func() {
		runtime.LockOSThread()
		err = iw.runGui()
	})
	return err
}

// todo: not working in game mode
func (iw *InputWindow) Show() {
	iw.window.SetFullscreen(0)
	iw.window.Hide()

	iw.window.Show()
	iw.window.Flash(sdl.FLASH_UNTIL_FOCUSED)
	iw.window.Raise()
	iw.window.SetFullscreen(sdl.WINDOW_FULLSCREEN_DESKTOP)
}

func (iw *InputWindow) runPendingEventsHandler() {
	for {
		iw.pendingEventsLoopLock.Lock()
		err := iw.handlePendingEventsBatch()
		if err != nil {
			log.WithError(err).Error("failed to handle events")
			continue
		}
	}
}

func (iw *InputWindow) runGui() error {
	var err error
	sdl.Do(func() {
		iw.window, err = sdl.CreateWindow(
			"DeckJoy Input",
			sdl.WINDOWPOS_UNDEFINED,
			sdl.WINDOWPOS_UNDEFINED,
			100,
			100,
			sdl.WINDOW_FULLSCREEN_DESKTOP|sdl.WINDOW_OPENGL,
		)
	})
	if err != nil {
		return err
	}
	defer func() {
		sdl.Do(func() {
			iw.window.Destroy()
		})
	}()

	sdl.Do(func() {
		iw.renderer, err = sdl.CreateRenderer(iw.window, -1, sdl.RENDERER_ACCELERATED)
	})
	if err != nil {
		return err
	}
	defer func() {
		sdl.Do(func() {
			iw.renderer.Destroy()
		})
	}()

	iw.keyboardGui = &KeyboardGUI{
		renderer:     iw.renderer,
		pixelPerUnit: 68,
		keyboard:     KeyboardAnsiTKL,
	}

	for {
		err = iw.loop()
		if errors.Is(err, QuitErr) {
			break
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (iw *InputWindow) loop() error {
	var err error
	sdl.Do(func() {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			err = iw.handleEvent(event)
			if err != nil {
				break
			}
		}
	})
	if err != nil {
		return err
	}
	if !iw.pendingEventsLoopLock.TryLock() {
		iw.pendingEventsLoopLock.Unlock()
	}

	sdl.Do(func() {
		err = iw.renderer.Clear()
	})
	if err != nil {
		return err
	}

	sdl.Do(func() {
		err = iw.renderer.SetDrawColor(0, 0, 0, 0x20)
	})
	if err != nil {
		return err
	}
	sdl.Do(func() {
		w, h := iw.window.GetSize()
		err = iw.renderer.FillRect(&sdl.Rect{W: w, H: h})
	})
	if err != nil {
		return err
	}

	err = iw.keyboardGui.Render(10, 390)
	if err != nil {
		return err
	}

	sdl.Do(func() {
		iw.renderer.Present()
		sdl.Delay(1)
	})
	return nil
}

func (iw *InputWindow) getTouchCords(event *sdl.TouchFingerEvent) (int32, int32) {
	windowW, windowH := iw.window.GetSize()
	x := int32(event.X * float32(windowW))
	y := int32(event.Y * float32(windowH))
	return x, y
}

func (iw *InputWindow) handlePendingEventsBatch() error {
	// filtering events
	// only keep the last TouchFingerEvent per finger, sdl detects double touches
	eventsToApply := []sdl.Event{}
	fingerEvents := map[sdl.FingerID]sdl.Event{}
	iw.pendingEventsBatchMutex.Lock()
	for _, event := range iw.pendingEventsBatch {
		switch t := event.(type) {
		case *sdl.TouchFingerEvent:
			fingerEvents[t.FingerID] = event
		default:
			eventsToApply = append(eventsToApply, event)
		}
	}
	iw.pendingEventsBatch = []sdl.Event{}
	iw.pendingEventsBatchMutex.Unlock()

	for _, event := range fingerEvents {
		eventsToApply = append(eventsToApply, event)
	}

	for _, event := range eventsToApply {
		switch t := event.(type) {
		case *sdl.MouseMotionEvent:
			err := iw.handleMouseMotionEvent(t)
			if err != nil {
				return err
			}
		case *sdl.MouseButtonEvent:
			err := iw.handleMouseButtonEvent(t)
			if err != nil {
				return err
			}
		case *sdl.TouchFingerEvent:
			err := iw.handleTouchFingerEvent(t)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (iw *InputWindow) storeEvent(event sdl.Event) {
	iw.pendingEventsBatchMutex.Lock()
	defer iw.pendingEventsBatchMutex.Unlock()
	iw.pendingEventsBatch = append(iw.pendingEventsBatch, event)
}

func (iw *InputWindow) handleEvent(event sdl.Event) error {
	switch t := event.(type) {
	case *sdl.QuitEvent:
		return QuitErr
	case *sdl.WindowEvent:
		if t.Event == sdl.WINDOWEVENT_FOCUS_LOST {
			log.Infof("input window lost focus")
			sdl.SetRelativeMouseMode(false)
			if err := steamworks.ActivateActionSetForDefaultController(steamworks.DefaultActionSet); err != nil {
				log.WithError(err).Error("failed to set DefaultActionSet")
			}
		} else if t.Event == sdl.WINDOWEVENT_FOCUS_GAINED {
			log.Infof("input window gained focus")
			sdl.SetRelativeMouseMode(true)
			if err := steamworks.ActivateActionSetForDefaultController(steamworks.TouchActionSet); err != nil {
				log.WithError(err).Error("failed to set TouchActionSet")
			}
		}
	case *sdl.MouseMotionEvent:
		copyEvent := *t
		iw.storeEvent(&copyEvent)
	case *sdl.MouseButtonEvent:
		copyEvent := *t
		iw.storeEvent(&copyEvent)
	case *sdl.TouchFingerEvent:
		copyEvent := *t
		iw.storeEvent(&copyEvent)
	}
	return nil
}

func (iw *InputWindow) handleMouseButtonEvent(t *sdl.MouseButtonEvent) error {
	if iw.deck.Mouse == nil {
		return nil
	}
	if t.Which > 100 {
		// ignore touchscreen "mouse" presses
		return nil
	}
	if t.Button == sdl.BUTTON_LEFT {
		if t.Type == sdl.MOUSEBUTTONDOWN {
			err := iw.deck.Mouse.PressButton(hid.MouseButtonLeft)
			if err != nil {
				return err
			}
		} else {
			err := iw.deck.Mouse.ReleaseButton(hid.MouseButtonLeft)
			if err != nil {
				return err
			}
		}
	}
	if t.Button == sdl.BUTTON_MIDDLE {
		if t.Type == sdl.MOUSEBUTTONDOWN {
			err := iw.deck.Mouse.PressButton(hid.MouseButtonMiddle)
			if err != nil {
				return err
			}
		} else {
			err := iw.deck.Mouse.ReleaseButton(hid.MouseButtonMiddle)
			if err != nil {
				return err
			}
		}
	}
	if t.Button == sdl.BUTTON_RIGHT {
		if t.Type == sdl.MOUSEBUTTONDOWN {
			err := iw.deck.Mouse.PressButton(hid.MouseButtonRight)
			if err != nil {
				return err
			}
		} else {
			err := iw.deck.Mouse.ReleaseButton(hid.MouseButtonRight)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func (iw *InputWindow) handleTouchFingerEvent(t *sdl.TouchFingerEvent) error {
	if iw.deck.Keyboard == nil {
		return nil
	}
	if t.Type == sdl.FINGERDOWN {
		x, y := iw.getTouchCords(t)
		key, err := iw.keyboardGui.GetKeyAt(x, y)
		if err != nil {
			// not a key, ignore press
			return nil
		}
		key.pressed = true
		iw.touches[t.FingerID] = key
		if key.Key != 0 {
			err = iw.deck.Keyboard.Press(key.Key)
			if err != nil {
				return err
			}
		}
		if key.ModKey != 0 {
			err = iw.deck.Keyboard.PressMod(key.ModKey)
			if err != nil {
				return err
			}
		}
	} else if t.Type == sdl.FINGERUP {
		key, found := iw.touches[t.FingerID]
		if found {
			key.pressed = false
			delete(iw.touches, t.FingerID)
			if key.Key != 0 {
				err := iw.deck.Keyboard.Release(key.Key)
				if err != nil {
					return err
				}
			}
			if key.ModKey != 0 {
				err := iw.deck.Keyboard.ReleaseMod(key.ModKey)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (iw *InputWindow) handleMouseMotionEvent(t *sdl.MouseMotionEvent) error {
	if iw.deck.Mouse == nil {
		return nil
	}

	err := iw.deck.Mouse.Move(int16(t.XRel), int16(t.YRel))
	if err != nil {
		return err
	}

	return nil
}
