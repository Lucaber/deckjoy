package setup

import (
	"fmt"
	"github.com/lucaber/deckjoy/pkg/deck"
	"github.com/lucaber/deckjoy/pkg/hid"
	"github.com/lucaber/deckjoy/pkg/usbgadget"
)

type AfterEnableHookFunc func() error

type Deck struct {
	gadget           *usbgadget.Gadget
	conf             *usbgadget.Config
	afterEnableHooks map[string]AfterEnableHookFunc
}

func NewDeck() (*Deck, error) {
	return &Deck{
		afterEnableHooks: map[string]AfterEnableHookFunc{},
	}, nil
}

func (d *Deck) SetupModules() error {
	return Modprobe("libcomposite")
}

func (d *Deck) setupGadget() error {
	gadget, err := usbgadget.CreateGadget("/sys/kernel/config/", "g.1")
	if err != nil {
		return err
	}
	d.gadget = gadget
	return nil
}

func (d *Deck) SetupGadget() error {
	err := d.setupGadget()
	if err != nil {
		return err
	}

	err = d.gadget.SetAttributes(usbgadget.GadgetAttributes{
		IDVendor:  0x1d6b,
		IDProduct: 0x0104,
		BCDDevice: 0x0100,
		BCDUSB:    0x0200,
	})

	serial, err := deck.SerialNumber()
	if err != nil || serial == "" {
		serial = "1"
	}

	err = d.gadget.SetStrings(usbgadget.GadgetStrings{
		Manufacturer: "Valve",
		Product:      "Steam Deck (DeckJoy)",
		SerialNumber: serial,
	})
	if err != nil {
		return err
	}

	d.conf, err = d.gadget.CreateConfig("c.1")
	if err != nil {
		return err
	}

	err = d.conf.SetAttributes(usbgadget.ConfigAttributes{
		MaxPower: 1000,
	})
	if err != nil {
		return err
	}
	err = d.conf.SetStrings(usbgadget.ConfigStrings{
		Configuration: "Steam Deck (DeckJoy)",
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *Deck) enable() error {
	if d.gadget == nil {
		return fmt.Errorf("gadget not setup")
	}

	// disable first
	_ = d.gadget.Disable()

	if err := d.gadget.Enable(); err != nil {
		return err
	}

	for _, hook := range d.afterEnableHooks {
		if err := hook(); err != nil {
			return err
		}
	}

	return nil
}

func (d *Deck) SetupJoystick(userPermissions bool) (string, error) {
	return d.SetupHidDevice("hid.joystick", hid.JoystickReportDesc, userPermissions)
}

func (d *Deck) SetupKeyboard(userPermissions bool) (string, error) {
	return d.SetupHidDevice("hid.keyboard", hid.KeyboardReportDesc, userPermissions)
}

func (d *Deck) SetupMouse(userPermissions bool) (string, error) {
	return d.SetupHidDevice("hid.mouse", hid.MouseReportDesc, userPermissions)
}

func (d *Deck) SetupHidDevice(name string, reportDesc []byte, userPermissions bool) (string, error) {
	if d.conf == nil {
		return "", fmt.Errorf("gadget not setup")
	}

	hid, err := d.gadget.CreateFunctionHID(name)
	if err != nil {
		return "", err
	}

	err = hid.SetAttributes(usbgadget.FunctionHIDAttributes{
		Protocol:     1,
		Subclass:     1,
		ReportLength: 8,
		ReportDesc:   reportDesc,
	})
	if err != nil {
		return "", err
	}

	// disable first, before adding new function to the config
	_ = d.gadget.Disable()

	err = d.conf.AddFunction(hid.FunctionGeneric)
	if err != nil {
		return "", err
	}

	if userPermissions {
		d.afterEnableHooks[name] = func() error {
			err = hid.SetDevicePermissions()
			if err != nil {
				return err
			}
			return nil
		}
	}

	err = d.enable()
	if err != nil {
		return "", fmt.Errorf("enable error %w", err)
	}

	path, err := hid.GetDevicePath()
	if err != nil {
		return "", err
	}

	return path, nil
}

func (d *Deck) Destroy() error {
	if d.gadget == nil {
		err := d.setupGadget()
		if err != nil {
			return err
		}
	}
	return d.gadget.Destroy(true)
}
