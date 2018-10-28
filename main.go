package main

import "github.com/faiface/pixel/pixelgl"

var stateCh = make(chan State)
var exitCh = make(chan bool)

type State int

const (
	Menu State = iota
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
			case Host:
				pixelgl.Run(runHost)
			case Game:
				pixelgl.Run(runGame)
			}
		default:
			pixelgl.Run(runMenu)
		}
	}
}
