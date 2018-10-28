package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	clients = make(map[int]*websocket.Conn)
)

func clientConnected(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("upgrade: %v\n", err)
		return
	}
	defer c.Close()

	cID := rand.Int()
	clients[cID] = c

	for {
		mt, rawMessage, err := c.ReadMessage()
		if err != nil {
			fmt.Printf("WSS Receive Error: %+v\n", err)
			break
		}

		var m = &message{}
		fmt.Printf("recv: %v\n", rawMessage)

		jsonErr := json.Unmarshal(rawMessage, m)
		if jsonErr != nil {
			fmt.Printf("WSS Json Unmarshal Error: %+v\n", jsonErr)
			break
		}

		fmt.Printf("Received message:%+v\n", m)

		// broadcast to players
		for i, p := range clients {
			err = p.WriteMessage(mt, rawMessage)
			if err != nil {
				fmt.Printf("WSS Send to Client[%d] Error: %+v\n", i, err)
				break
			}

			m.ActionSquare.set(m.BlockId)

		}
		delete(clients, cID)
	}
}

func runServer(lt net.Listener) {

	// @TODO http.HandleFunc("/", httpHandler)
	http.HandleFunc("/ws", clientConnected)

	server := http.Server{
		Addr: ":" + Config.Port,
	}
	fmt.Printf("DEBUG %+v\n", lt)

	fmt.Printf("Server is listening on %v\n", server.Addr)
	err := http.Serve(lt, nil)
	if err != nil {
		fmt.Printf("Serve: ", err)
	}
}
