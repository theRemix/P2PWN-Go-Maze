package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
)

type OpCode int

const (
	_ OpCode = iota
	SetActionSquare
	StateOld
	StateNew
)

var clients = make(map[int]*Client)

func clientConnected(w http.ResponseWriter, r *http.Request) {
	cID := rand.Int()
	newClient := &Client{
		ID:      cID,
		Updated: worldState.LastUpdated,
	}
	clients[cID] = newClient

	json.NewEncoder(w).Encode(newClient)
}

func clientUpdate(w http.ResponseWriter, r *http.Request) {
	var clientToUpdate = &Client{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(clientToUpdate); err != nil {
		fmt.Printf("Error reading client update: %s\n", err)
		return
	}

	updateResponse := &UpdateResponse{
		State: StateOld,
	}

	if clientToUpdate.NeedsUpdate() {
		updateResponse = &UpdateResponse{
			State:      StateNew,
			Updated:    worldState.LastUpdated,
			WorldTiles: worldState.Tiles,
		}
	}

	json.NewEncoder(w).Encode(updateResponse)
}

func clientActed(w http.ResponseWriter, r *http.Request) {
	var ca = &ClientAction{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(ca); err != nil {
		fmt.Printf("Error reading client action: %s", err)
		return
	}

	ca.ActionSquare.active = true
	ca.ActionSquare.set(ca.BlockId)
}

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, home)
}

func runServer(lt net.Listener) {
	client.ID = 0

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/srv/connect", clientConnected)
	http.HandleFunc("/srv/update", clientUpdate)
	http.HandleFunc("/srv/action", clientActed)
	http.HandleFunc("/srv/health", health)

	server := http.Server{
		Addr: ":" + Config.Port,
	}

	fmt.Printf("Server is listening on %v\n", P2pwn.EntryURL)
	err := server.Serve(lt)
	if err != nil {
		fmt.Printf("Error Serve:%v\n", err)
	}
}
