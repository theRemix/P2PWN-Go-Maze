package main

import (
	"flag"
	"github.com/faiface/pixel/pixelgl"
)

//go:generate go run includes/include.go

var stateCh = make(chan State)
var exitCh = make(chan bool)

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
	setConfig(&Config.P2pwn, "p2pwn", "https://p2pwithme.2018.nodeknockout.com", "P2PWN Service Address")

	flag.Parse()

	for {
		select {
		case <-exitCh:
			return
		case state := <-stateCh: // not a great way to do this! #hackathon
			switch state {
			case Menu:
				pixelgl.Run(runMenu)
			case Join:
				pixelgl.Run(runJoin)
			case Game:
				pixelgl.Run(runGame)
			case Host:
				pixelgl.Run(runHost)
			}
		default:
			pixelgl.Run(runMenu)
		}
	}
}
