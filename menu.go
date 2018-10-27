package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
)

const fontFace = "fonts/zorque.ttf"

func loadTTF(path string, size float64) (font.Face, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	font, err := truetype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(font, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
}

func drawButtons(win *pixelgl.Window, titleTxt, hostTxt, joinTxt *text.Text) {
	hostTxt.Clear()
	joinTxt.Clear()
	hostTxt.WriteString("Host New Room")
	joinTxt.WriteString("Join Room")
	titleTxt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(pixel.V(550, -50))))
	hostTxt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(pixel.V(670, 200))))
	joinTxt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(pixel.V(570, 300))))
}

func runMenu() {
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

	drawButtons(win, titleTxt, hostTxt, joinTxt)

	hostBounds := pixel.R(195, 285, 854, 336)
	joinBounds := pixel.R(297, 183, 730, 237)

	for !win.Closed() {
		if win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ) {
			return
		}

		if hostBounds.Contains(win.MousePosition()) {
			win.Clear(colornames.Firebrick)
			hostTxt.Color = colornames.Darkturquoise
			drawButtons(win, titleTxt, hostTxt, joinTxt)
		} else if joinBounds.Contains(win.MousePosition()) {
			win.Clear(colornames.Firebrick)
			joinTxt.Color = colornames.Darkturquoise
			drawButtons(win, titleTxt, hostTxt, joinTxt)
		} else if hostTxt.Color != colornames.Darkkhaki {
			win.Clear(colornames.Firebrick)
			hostTxt.Color = colornames.Darkkhaki
			drawButtons(win, titleTxt, hostTxt, joinTxt)
		} else if joinTxt.Color != colornames.Darkkhaki {
			win.Clear(colornames.Firebrick)
			joinTxt.Color = colornames.Darkkhaki
			drawButtons(win, titleTxt, hostTxt, joinTxt)
		}

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			if hostBounds.Contains(win.MousePosition()) {
				fmt.Printf("HOST %v", win.MousePosition())
			} else if joinBounds.Contains(win.MousePosition()) {
				fmt.Printf("JOIN %v", win.MousePosition())
			}
		}

		win.Update()
	}
}