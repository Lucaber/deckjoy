package gui

import (
	"bytes"
	"image"
	"image/color"

	"github.com/veandco/go-sdl2/sdl"
	"golang.org/x/image/bmp"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// renderTextSurface renders text into an SDL surface using basicfont.
// This avoids depending on sdl_ttf or sdl_image which are not installed on the Steam Deck by default.
func renderTextSurface(text string, w, h int, col color.RGBA, dot fixed.Point26_6) (*sdl.Surface, error) {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  dot,
	}
	d.DrawString(text)

	bmpBytes := bytes.Buffer{}
	if err := bmp.Encode(&bmpBytes, img); err != nil {
		return nil, err
	}

	rw, err := sdl.RWFromMem(bmpBytes.Bytes())
	if err != nil {
		return nil, err
	}

	return sdl.LoadBMPRW(rw, true)
}
