package main

import (
	"fmt"
	"net/url"

	"golang.org/x/net/websocket"
)

var clientCh = make(chan actionSquare)

func clientConnect(rawurl string) {
	origin := "http://localhost/"
	url, urlErr := url.Parse(rawurl)
	url.Scheme = "ws" // coerce into ws://

	if urlErr != nil {
		fmt.Printf("Error Parsing URL: %s\n", urlErr)
		return
	}

	ws, err := websocket.Dial(url.String(), "", origin)
	if err != nil {
		fmt.Printf("Error Dialing Host: %s\n", err)
		return
	}

	go func() {
		for {
			select {
			case as := <-clientCh:
				if err := websocket.JSON.Send(ws, as); err != nil {
					fmt.Printf("WSS Send to Server Error: %+v\n", err)
					break
				}
			}
		}
	}()

	for {
		var m = &message{}

		// receive a message using the codec
		if err := websocket.JSON.Receive(ws, m); err != nil {
			fmt.Printf("WSS Receive Error: %+v\n", err)
			break
		}

		fmt.Printf("Received message:%+v\n", m)
		// @TODO handle message on local game
	}

}
