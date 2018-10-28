package main

import "github.com/faiface/pixel/pixelgl"

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

func main() {
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
