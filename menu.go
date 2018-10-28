package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

func drawMenuButtons(win *pixelgl.Window, titleTxt, hostTxt, joinTxt *text.Text) {
	hostTxt.Clear()
	joinTxt.Clear()
	hostTxt.WriteString("Host New Room")
	joinTxt.WriteString("Join Room")
	titleTxt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(pixel.V(550, -50))))
	hostTxt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(pixel.V(670, 200))))
	joinTxt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(pixel.V(570, 300))))
}

func runMenu() {
	const fontFace = font1 // zorque.ttf

	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1024, 768),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)

	face, err := loadTTF(fontFace, 80)
	if err != nil {
		panic(err)
	}

	win.Clear(colornames.Firebrick)

	atlas := text.NewAtlas(face, text.ASCII)

	titleTxt := text.New(pixel.V(350, 100), atlas)
	titleTxt.Color = colornames.Lightgrey
	titleTxt.WriteString("Go Maze!")

	hostTxt := text.New(pixel.V(350, 100), atlas)
	hostTxt.Color = colornames.Darkkhaki

	joinTxt := text.New(pixel.V(350, 100), atlas)
	joinTxt.Color = colornames.Darkkhaki

	drawMenuButtons(win, titleTxt, hostTxt, joinTxt)

	hostBounds := pixel.R(195, 285, 854, 336)
	joinBounds := pixel.R(297, 183, 730, 237)

	for !win.Closed() {
		if win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ) {
			go func() { exitCh <- true }()
			return
		}

		if hostBounds.Contains(win.MousePosition()) {
			win.Clear(colornames.Firebrick)
			hostTxt.Color = colornames.Darkturquoise
			drawMenuButtons(win, titleTxt, hostTxt, joinTxt)
		} else if joinBounds.Contains(win.MousePosition()) {
			win.Clear(colornames.Firebrick)
			joinTxt.Color = colornames.Darkturquoise
			drawMenuButtons(win, titleTxt, hostTxt, joinTxt)
		} else if hostTxt.Color != colornames.Darkkhaki {
			win.Clear(colornames.Firebrick)
			hostTxt.Color = colornames.Darkkhaki
			drawMenuButtons(win, titleTxt, hostTxt, joinTxt)
		} else if joinTxt.Color != colornames.Darkkhaki {
			win.Clear(colornames.Firebrick)
			joinTxt.Color = colornames.Darkkhaki
			drawMenuButtons(win, titleTxt, hostTxt, joinTxt)
		}

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			if hostBounds.Contains(win.MousePosition()) {
				go func() { stateCh <- Host }()
				return
			} else if joinBounds.Contains(win.MousePosition()) {
				go func() { stateCh <- Join }()
				return
			}
		}

		win.Update()
	}
	win.Destroy()
}
