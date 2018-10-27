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

	atlas := text.NewAtlas(face, text.ASCII)
	txt := text.New(pixel.V(50, 500), atlas)

	txt.Color = colornames.Lightgrey

	txt.WriteString("Go Maze!")
	txt.WriteRune('\n')
	txt.WriteRune('\n')
	txt.WriteString("Host New Room")
	txt.WriteRune('\n')
	txt.WriteString("Join Room")

	win.Clear(colornames.Firebrick)
	txt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(txt.Bounds().Center())))

	hostBounds := pixel.R(183, 313, 843, 367)
	joinBounds := pixel.R(183, 233, 616, 287)

	for !win.Closed() {
		if win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ) {
			return
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
