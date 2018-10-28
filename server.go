package main

import (
	"fmt"
	// "math/rand"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	clients  = make(map[int]*websocket.Conn)
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// serveWs handles websocket requests from the peer.
func clientConnected(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

// cID := rand.Int()
// clients[cID] = c

// for {
// 	mt, rawMessage, err := c.ReadMessage()
// 	if err != nil {
// 		fmt.Printf("WSS Receive Error: %+v\n", err)
// 		break
// 	}

// 	var m = &message{}
// 	fmt.Printf("recv: %v\n", rawMessage)

// 	jsonErr := json.Unmarshal(rawMessage, m)
// 	if jsonErr != nil {
// 		fmt.Printf("WSS Json Unmarshal Error: %+v\n", jsonErr)
// 		break
// 	}

// 	fmt.Printf("Received message:%+v\n", m)

// 	// broadcast to players
// 	for i, p := range clients {
// 		err = p.WriteMessage(mt, rawMessage)
// 		if err != nil {
// 			fmt.Printf("WSS Send to Client[%d] Error: %+v\n", i, err)
// 			break
// 		}

// 		m.ActionSquare.set(m.BlockId)

// 	}
// 	delete(clients, cID)
// }

func runServer(lt net.Listener) {
	hub := newHub()
	go hub.run()
	// @TODO http.HandleFunc("/", httpHandler)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		clientConnected(hub, w, r)
	})
	err := http.ListenAndServe(":"+Config.Port, nil)
	if err != nil {
		fmt.Printf("ListenAndServe: %v", err)
	}

	// server := http.Server{
	// 	Addr: ":" + Config.Port,
	// }

	// fmt.Printf("Server is listening on %v\n", server.Addr)
	// server.Serve(lt)
}
