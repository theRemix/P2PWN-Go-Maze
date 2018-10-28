package main

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"

	"golang.org/x/net/websocket"
)

type message struct {
	Op           string
	ActionSquare *actionSquare
	BlockId      int
}

var clients = make(map[int]*websocket.Conn)

func clientConnected(ws *websocket.Conn) {
	cID := rand.Int()
	clients[cID] = ws

	for {
		var m = &message{}

		// receive a message using the codec
		if err := websocket.JSON.Receive(ws, m); err != nil {
			fmt.Printf("WSS Receive Error: %+v\n", err)
			break
		}

		fmt.Printf("Received message:%+v\n", m)

		// broadcast to players
		for i, c := range clients {
			if err := websocket.JSON.Send(c, m); err != nil {
				fmt.Printf("WSS Send to Client[%d] Error: %+v\n", i, err)
				break
			}
		}

		// @TODO handle message on local game
	}

	delete(clients, cID)
}

func runServer(lt net.Listener) {

	http.Handle("/", websocket.Handler(clientConnected))

	server := http.Server{
		Addr: ":" + Config.Port,
	}

	fmt.Printf("Server is listening on %v\n", server.Addr)
	server.Serve(lt)
}
