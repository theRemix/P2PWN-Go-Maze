package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

type p2pHost struct {
	ID          string `json:"id"`
	AppName     string `json:"app_name"`
	DisplayName string `json:"display_name"`
	EntryURL    string `json:"entry_url"`
}

func (h *p2pHost) getRect(idx int) pixel.Rect {
	return pixel.R(100, float64(470+(idx*-90)), 800, float64(440+((idx-1)*-90)))
}

type p2pHostList []*p2pHost

func drawHostList(win *pixelgl.Window, atlas *text.Atlas, hosts *p2pHostList) {
	for idx, host := range *hosts {

		hostTxt := text.New(pixel.V(0, 0), atlas)
		hostTxt.Color = colornames.Darkkhaki
		hostTxt.WriteString(host.DisplayName)
		r := host.getRect(idx)
		hostTxt.Draw(win, pixel.IM.Moved(pixel.V(r.Min.X, r.Min.Y)))

	}
}

func getHostList() (*p2pHostList, error) {
	hostList := &p2pHostList{}
	hostRes, err := http.Get(Config.P2pwn + "/api/hosts")
	if err != nil {
		fmt.Printf("Error Connecting to P2PWN Service: %v\n", err)
	}

	defer hostRes.Body.Close()
	if err := json.NewDecoder(hostRes.Body).Decode(hostList); err != nil {
		fmt.Println("Unmarshal P2PWN Response Error:", err)
	}

	return hostList, err
}

func runJoin() {
	const titleFont = "fonts/zorque.ttf"
	const hostFont = "fonts/gomarice_game_continue_02.ttf"

	cfg := pixelgl.WindowConfig{
		Title:  "Join Game",
		Bounds: pixel.R(0, 0, 1024, 768),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)

	titleFontFace, titleFontErr := loadTTF(titleFont, 80)
	if titleFontErr != nil {
		panic(titleFontErr)
	}

	hostFontFace, hostFontErr := loadTTF(hostFont, 80)
	if hostFontErr != nil {
		panic(hostFontErr)
	}

	win.Clear(colornames.Firebrick)
	titleAtlas := text.NewAtlas(titleFontFace, text.ASCII)
	atlas := text.NewAtlas(hostFontFace, text.ASCII)

	titleTxt := text.New(pixel.V(350, 100), titleAtlas)
	titleTxt.Color = colornames.Lightgrey
	titleTxt.WriteString("Join Game")
	titleTxt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(pixel.V(560, -100))))

	backTxt := text.New(pixel.V(150, 30), titleAtlas)
	backTxt.Color = colornames.Lightgrey
	backTxt.WriteString("Back")
	backTxt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(pixel.V(610, -270))))

	hosts, getHostsErr := getHostList()
	if getHostsErr != nil {
		fmt.Printf("Problem getting list of game hosts. %v\n", getHostsErr)
	}

	drawHostList(win, atlas, hosts)

	for !win.Closed() {
		if win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ) {
			go func() { exitCh <- true }()
			return
		}

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			if pixel.R(0, 650, 300, 800).Contains(win.MousePosition()) {
				go func() { stateCh <- Menu }()
				return
			}

			for idx, host := range *hosts {
				if host.getRect(idx).Contains(win.MousePosition()) {
					go func() { stateCh <- Game }()
					go clientConnect(host.EntryURL)
					return
				}
			}
		}

		win.Update()
	}
	win.Destroy()

}
