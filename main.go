package main

import (
	"flag"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

//go:generate go run includes/include.go

// CHANGE ME
const (
	appName    = "p2pwn-go-maze"
	appRelease = "DEVELOPMENT"
)

var (
	stateCh = make(chan State)
	exitCh  = make(chan bool)
	state   = Menu
)

type State int

const (
	_ State = iota
	Menu
	Join
	Host
	Game
)

// App  Config
type appConfig struct {
	AppName        string `json:"app_name"`        // for grouping rooms in P2PWN
	DisplayName    string `json:"display_name"`    // used to display in P2PWN lobby
	Release        string `json:"release"`         // "PRODUCTION", "DEVELOPMENT"
	EntryURL       string `json:"entry_url"`       // url used as the entrypoint for your app, supplied by localtunnel
	HealthCheckURL string `json:"healthcheck_url"` // health endpoint
	Port           string // Server Listening Port
	P2pwn          string // P2PWN Service Address
}

var Config = &appConfig{}

func main() {
	setConfig(&Config.AppName, "name", appName, "Name of this app")
	setConfig(&Config.Port, "port", "3000", "Port for server to listen on")
	setConfig(&Config.P2pwn, "p2pwn", "https://p2pwn-production.herokuapp.com", "P2PWN Service Address")

	flag.Parse()

	cfg := pixelgl.WindowConfig{
		Title: "Go Maze",
		// Bounds: pixel.R(0, 0, 1024, 768),
		Bounds: pixel.R(0, 0, float64(width)*scale, float64(height)*scale),
		VSync:  true,
	}

	pixelgl.Run(func() {
		win, err := pixelgl.NewWindow(cfg)
		if err != nil {
			panic(err)
		}
		win.SetSmooth(true)

		go func() {
			for {
				select {
				case <-exitCh:
					win.Destroy()
					return
				case state = <-stateCh: // not a great way to do this! #hackathon
				}
			}
		}()

		for !win.Closed() {

			if win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ) {
				return
			}

			switch state {
			case Menu:
				runMenu(win)
			case Join:
				runJoin(win)
			case Game:
				runGame(win)
			case Host:
				runHost(win)
			default:
				runMenu(win)
			}
		}

		win.Destroy()

	})

}
