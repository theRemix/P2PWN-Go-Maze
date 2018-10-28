package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"time"
)

var clients = make(map[int]*Client)

func clientConnected(w http.ResponseWriter, r *http.Request) {
	cID := rand.Int()
	newClient := &Client{
		ID:      cID,
		Updated: time.Now(),
	}
	clients[cID] = newClient

	json.NewEncoder(w).Encode(newClient)
}

func clientUpdate(w http.ResponseWriter, r *http.Request) {
	var client = &Client{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(client); err != nil {
		fmt.Printf("Error reading ping body: %s", err)
		return
	}

	clients[client.ID].Updated = time.Now()

	json.NewEncoder(w).Encode(world)
}

func clientActed(w http.ResponseWriter, r *http.Request) {
	var ca = &ClientAction{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(ca); err != nil {
		fmt.Printf("Error reading client action: %s", err)
		return
	}

	world[ca.ActionSquare.X][ca.ActionSquare.Y] = ca.BlockId
}

func runServer(lt net.Listener) {
	client.ID = 0

	// @TODO http.HandleFunc("/", homeHandler)
	http.HandleFunc("/srv/connect", clientConnected)
	http.HandleFunc("/srv/update", clientUpdate)
	http.HandleFunc("/srv/action", clientActed)

	server := http.Server{
		Addr: ":" + Config.Port,
	}

	fmt.Printf("Server is listening on %v\n", server.Addr)
	err := server.Serve(lt)
	if err != nil {
		fmt.Printf("Error Serve:%v\n", err)
	}
}
