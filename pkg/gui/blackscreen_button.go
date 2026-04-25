package gui

import (
	"image/color"

	"github.com/veandco/go-sdl2/sdl"
	"golang.org/x/image/math/fixed"
)

func (iw *InputWindow) preRenderBlackScreenButton() error {
	iw.blackScreenBtnRect = sdl.Rect{
		X: 10,
		Y: 20,
		W: 180,
		H: 50,
	}

	surface, err := renderTextSurface("Black Screen", 180, 50, color.RGBA{255, 255, 255, 255}, fixed.Point26_6{X: fixed.I(30), Y: fixed.I(30)})
	if err != nil {
		return err
	}
	iw.blackScreenBtnSurface = surface

	return nil
}

func (iw *InputWindow) renderBlackScreenButton() error {
	var err error
	sdl.Do(func() {
		if err = iw.renderer.SetDrawColor(0xFF, 0xFF, 0xFF, 0xFF); err != nil {
			return
		}
		if err = iw.renderer.DrawRect(&iw.blackScreenBtnRect); err != nil {
			return
		}

		texture, texErr := iw.renderer.CreateTextureFromSurface(iw.blackScreenBtnSurface)
		if texErr != nil {
			err = texErr
			return
		}
		defer texture.Destroy()

		err = iw.renderer.Copy(texture, &sdl.Rect{W: iw.blackScreenBtnRect.W, H: iw.blackScreenBtnRect.H}, &iw.blackScreenBtnRect)
	})
	return err
}
