package gui

import (
	"bytes"
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"golang.org/x/image/bmp"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
)

type KeyboardGUI struct {
	renderer     *sdl.Renderer
	pixelPerUnit float64
	keyboard     Keyboard
}

func (kr *KeyboardGUI) RenderKey(key *Key) error {
	if key.Hidden {
		return nil
	}

	if key.pressed {
		if err := kr.renderer.SetDrawColor(0x00, 0xff, 0xff, 0xff); err != nil {
			return err
		}
	} else {
		if err := kr.renderer.SetDrawColor(0xff, 0x00, 0x00, 0xff); err != nil {
			return err
		}
	}

	keyRect := &sdl.Rect{
		X: key.renderX,
		Y: key.renderY,
		W: key.renderW,
		H: key.renderH,
	}

	if err := kr.renderer.DrawRect(keyRect); err != nil {
		return err
	}

	if key.renderSurface != nil {
		texture, err := kr.renderer.CreateTextureFromSurface(key.renderSurface)
		if err != nil {
			return err
		}

		if err := kr.renderer.Copy(texture, &sdl.Rect{W: keyRect.W, H: keyRect.H}, keyRect); err != nil {
			return err
		}

		if err := texture.Destroy(); err != nil {
			return err
		}
	}
	return nil
}

func (kr *KeyboardGUI) PreRender(paddingX, paddingY int) error {
	y := float64(paddingY) / kr.pixelPerUnit
	for rowNum := range kr.keyboard.Rows {
		x := float64(paddingX) / kr.pixelPerUnit
		for keyNum := range kr.keyboard.Rows[rowNum].Keys {
			key := &kr.keyboard.Rows[rowNum].Keys[keyNum]

			key.renderX = int32(kr.pixelPerUnit * x)
			key.renderY = int32(kr.pixelPerUnit * y)
			key.renderW = int32(kr.pixelPerUnit * key.WidthUnits)
			key.renderH = int32(kr.pixelPerUnit)

			if key.Text != "" {
				// rendering text without sdl_ttf or sdl_image
				// not installed on steam deck by default
				img := image.NewRGBA(image.Rect(0, 0, int(key.renderW), int(key.renderH)))
				location := fixed.Point26_6{fixed.I(20), fixed.I(30)}
				d := &font.Drawer{
					Dst:  img,
					Src:  image.NewUniform(color.RGBA{200, 100, 0, 255}),
					Face: basicfont.Face7x13,
					Dot:  location,
				}

				d.DrawString(key.Text)

				bmpBytes := bytes.Buffer{}
				err := bmp.Encode(&bmpBytes, img)
				if err != nil {
					return err
				}

				rw, err := sdl.RWFromMem(bmpBytes.Bytes())
				if err != nil {
					return err
				}
				surface, err := sdl.LoadBMPRW(rw, true)
				if err != nil {
					return err
				}
				key.renderSurface = surface
			}

			x += key.WidthUnits
		}
		y += 1
	}
	return nil
}

func (kr *KeyboardGUI) Render(paddingX, paddingY int) error {
	if !kr.keyboard.PreRendered {
		if err := kr.PreRender(paddingX, paddingY); err != nil {
			return err
		}
		kr.keyboard.PreRendered = true
	}

	for rowNum := range kr.keyboard.Rows {
		for keyNum := range kr.keyboard.Rows[rowNum].Keys {
			key := &kr.keyboard.Rows[rowNum].Keys[keyNum]
			var err error
			sdl.Do(func() {
				err = kr.RenderKey(key)
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

var KeyNotFoundErr = fmt.Errorf("key not found")

func (kr *KeyboardGUI) GetKeyAt(x int32, y int32) (*Key, error) {
	for rowNum := range kr.keyboard.Rows {
		for keyNum := range kr.keyboard.Rows[rowNum].Keys {
			key := &kr.keyboard.Rows[rowNum].Keys[keyNum]
			if x > key.renderX && x < (key.renderX+key.renderW) &&
				y > key.renderY && y < (key.renderY+key.renderH) {
				return key, nil
			}
		}
	}
	return nil, KeyNotFoundErr
}
