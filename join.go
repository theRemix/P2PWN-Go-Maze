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
	return pixel.R(100, float64(440+(idx*-40)), 800, float64(440+((idx-1)*-40)))
}

type p2pHostList []*p2pHost

func drawHostList(win *pixelgl.Window, atlas *text.Atlas, hosts *p2pHostList) {
	for idx, host := range *hosts {

		hostTxt := text.New(pixel.V(0, 0), atlas)
		hostTxt.Color = colornames.Darkkhaki
		hostTxt.WriteString(host.DisplayName)
		r := host.getRect(len(*hosts) - 1 - idx)
		hostTxt.Draw(win, pixel.IM.Moved(pixel.V(r.Min.X, r.Min.Y)))

	}
}

func getHostList() (*p2pHostList, error) {
	hostList := &p2pHostList{}
	hostRes, err := http.Get(Config.P2pwn + "/api/hosts?app_name=" + Config.AppName)
	if err != nil {
		fmt.Printf("Error Connecting to P2PWN Service: %v\n", err)
	}

	defer hostRes.Body.Close()
	if err := json.NewDecoder(hostRes.Body).Decode(hostList); err != nil {
		fmt.Println("Unmarshal P2PWN Response Error:", err)
	}

	return hostList, err
}

func runJoin(win *pixelgl.Window) {
	const titleFont = font1 // zorque.ttf
	const hostFont = font2  // gomarice_game_continue_02.ttf

	titleFontFace, titleFontErr := loadTTF(titleFont, 80)
	if titleFontErr != nil {
		panic(titleFontErr)
	}

	hostFontFace, hostFontErr := loadTTF(hostFont, 40)
	if hostFontErr != nil {
		panic(hostFontErr)
	}

	win.Clear(colornames.Firebrick)
	titleAtlas := text.NewAtlas(titleFontFace, text.ASCII)
	atlas := text.NewAtlas(hostFontFace, text.ASCII)

	titleTxt := text.New(pixel.V(350, 100), titleAtlas)
	titleTxt.Color = colornames.Lightgrey
	titleTxt.WriteString("Join Game")
	titleTxt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(pixel.V(800, -110))))

	backTxt := text.New(pixel.V(350, 100), titleAtlas)
	backTxt.Color = colornames.Lightgrey
	backTxt.WriteString("Back")
	backTxt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(pixel.V(120, -110))))

	hosts, getHostsErr := getHostList()
	if getHostsErr != nil {
		fmt.Printf("Problem getting list of game hosts. %v\n", getHostsErr)
		win.Update()
		go func() { stateCh <- Menu }()
		return
	}

	drawHostList(win, atlas, hosts)

	for state == Join {
		if win.Closed() || win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ) {
			go func() { exitCh <- true }()
			return
		}

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			if pixel.R(680, 470, 960, 600).Contains(win.MousePosition()) {
				win.Update()
				go func() { stateCh <- Menu }()
				return
			}

			for idx, host := range *hosts {
				if host.getRect(len(*hosts) - 1 - idx).Contains(win.MousePosition()) {
					go func() { stateCh <- Game }()
					go clientConnect(host.EntryURL)
					return
				}
			}
		}

		win.Update()
	}
}
