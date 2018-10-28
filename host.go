package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/localtunnel/go-localtunnel"
	"golang.org/x/image/colornames"
)

// CHANGE ME
const appName = "p2pwn-go-maze"
const appRelease = "DEVELOPMENT"

// P2PWN Service Config
var P2pwn = &p2pwnConfig{}

type p2pwnConfig struct { // all values will be provided by P2PWN
	ID             string `json:"id"`              // public id assigned by P2PWN service
	AccessToken    string `json:"access_token"`    // private access token needed to perform actions on this host
	AppName        string `json:"app_name"`        // for grouping rooms in P2PWN
	DisplayName    string `json:"display_name"`    // used to display in P2PWN lobby
	EntryURL       string `json:"entry_url"`       // url used as the entrypoint for your app, supplied by localtunnel
	HealthCheckURL string `json:"healthcheck_url"` // health endpoint
}

func runHost(win *pixelgl.Window) {
	const uiFont = font1    // zorque.ttf
	const inputFont = font2 // gomarice_game_continue_02.ttf

	uiFontFace, uiFontFaceErr := loadTTF(uiFont, 60)
	if uiFontFaceErr != nil {
		panic(uiFontFaceErr)
	}

	inputFontFace, inputFontFaceErr := loadTTF(inputFont, 50)
	if inputFontFaceErr != nil {
		panic(inputFontFaceErr)
	}

	win.Clear(colornames.Firebrick)
	uiAtlas := text.NewAtlas(uiFontFace, text.ASCII)
	inputAtlas := text.NewAtlas(inputFontFace, text.ASCII)

	titleTxt := text.New(pixel.V(350, 100), uiAtlas)
	titleTxt.Color = colornames.Lightgrey
	titleTxt.WriteString("Host Game")

	labelTxt := text.New(pixel.V(350, 100), uiAtlas)
	labelTxt.Color = colornames.Lightgrey
	labelTxt.WriteString("Enter Server Name:")

	serverNameTxt := text.New(pixel.V(350, 100), inputAtlas)
	serverNameTxt.Color = colornames.Darkkhaki

	fps := time.Tick(time.Second / 120)

	serverName := []string{}

	for !win.Closed() {
		key := win.Typed()
		if win.JustPressed(pixelgl.KeyEnter) || win.Repeated(pixelgl.KeyEnter) {
			break
		} else if win.JustPressed(pixelgl.KeyBackspace) || win.Repeated(pixelgl.KeyBackspace) {
			serverName = serverName[:len(serverName)-1]
		} else if key != "" {
			serverName = append(serverName, key)
		}

		serverNameTxt.Clear()
		serverNameTxt.WriteString(strings.Join(serverName[:], ""))

		win.Clear(colornames.Midnightblue)
		titleTxt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(titleTxt.Bounds().Center().Sub(pixel.V(0, 200)))))
		labelTxt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(labelTxt.Bounds().Center().Sub(pixel.V(0, 100)))))
		serverNameTxt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(serverNameTxt.Bounds().Center().Add(pixel.V(0, 70)))))
		win.Update()

		<-fps
	}

	Config.DisplayName = strings.Join(serverName[:], "")
	Config.Release = appRelease

	port, portErr := strconv.Atoi(Config.Port)
	if portErr != nil {
		fmt.Printf("Invalid Port config: %s -> %v \n", Config.Port, port)
		os.Exit(1)
		return
	}

	lt, ltErr := localtunnel.Listen(localtunnel.Options{
		// 	Subdomain: Config.AppName,
	})
	if ltErr != nil {
		fmt.Printf("Error creating localtunnel: %v\n", ltErr)
		os.Exit(1)
		return
	}

	Config.EntryURL = lt.URL()
	Config.HealthCheckURL = lt.URL() + "/srv/health"
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

	// fmt.Printf("DEBUG P2PWN : %+v\n", P2pwn)

	labelTxt.Clear()
	labelTxt.Color = colornames.Darkcyan
	labelTxt.WriteString("Connected!")
	labelTxt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(pixel.V(630, 300))))
	win.Update()

	go func() { stateCh <- Game }()
	go runServer(lt)
}
