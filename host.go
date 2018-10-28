package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/localtunnel/go-localtunnel"
	"golang.org/x/image/colornames"
)

// CHANGE ME
const appName = "p2pwn-go-maze"
const displayName = "Go Maze"
const appRelease = "DEVELOPMENT"

// P2PWN Service Config
var P2pwn = &p2pwnConfig{}

type p2pwnConfig struct { // all values will be provided by P2PWN
	ID          string `json:"id"`           // public id assigned by P2PWN service
	AccessToken string `json:"access_token"` // private access token needed to perform actions on this host

	AppName     string `json:"app_name"`     // for grouping rooms in P2PWN
	DisplayName string `json:"display_name"` // used to display in P2PWN lobby
	EntryURL    string `json:"entry_url"`    // url used as the entrypoint for your app, supplied by localtunnel
}

func runHost() {
	const font = "fonts/zorque.ttf"

	cfg := pixelgl.WindowConfig{
		Title:  "Host Game",
		Bounds: pixel.R(0, 0, 1024, 768),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)

	fontFace, fontFaceErr := loadTTF(font, 80)
	if fontFaceErr != nil {
		panic(fontFaceErr)
	}

	win.Clear(colornames.Firebrick)
	atlas := text.NewAtlas(fontFace, text.ASCII)

	titleTxt := text.New(pixel.V(350, 100), atlas)
	titleTxt.Color = colornames.Lightgrey
	titleTxt.WriteString("Host Game")
	titleTxt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(pixel.V(600, -100))))

	statusTxt := text.New(pixel.V(350, 100), atlas)
	statusTxt.Color = colornames.Darkkhaki
	statusTxt.WriteString("Creating Server")
	statusTxt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(pixel.V(730, 200))))
	win.Update()

	Config.DisplayName = displayName
	Config.Release = appRelease

	port, portErr := strconv.Atoi(Config.Port)
	if portErr != nil {
		fmt.Printf("Invalid Port config: %s -> %v \n", Config.Port, port)
		os.Exit(1)
		return
	}

	lt, ltErr := localtunnel.Listen(localtunnel.Options{
		Subdomain: Config.AppName,
	})
	if ltErr != nil {
		fmt.Printf("Error creating localtunnel: %v\n", ltErr)
		os.Exit(1)
		return
	}

	Config.EntryURL = lt.URL()
	fmt.Printf("Connected to LT: %v\n", lt.URL())

	payload, _ := json.Marshal(Config)

	p2pwnRes, p2pwnErr := http.Post(Config.P2pwn+"/api/connect", "application/json", bytes.NewBuffer(payload))
	if p2pwnErr != nil {
		fmt.Printf("Error Connecting to P2PWN Service: %v\n", p2pwnErr)
		os.Exit(1)
		return
	}

	defer p2pwnRes.Body.Close()
	if err := json.NewDecoder(p2pwnRes.Body).Decode(P2pwn); err != nil {
		fmt.Println("Unmarshal P2PWN Response Error:", err)
		os.Exit(1)
		return
	}

	fmt.Printf("P2PWN is Ready: %+v\n", P2pwn)

	statusTxt.Clear()
	statusTxt.Color = colornames.Darkcyan
	statusTxt.WriteString("Connected!")
	statusTxt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(pixel.V(630, 300))))
	win.Update()
	win.Destroy()

	go func() { stateCh <- Game }()
	go runServer(lt)
}
