package main

import (
	"fmt"
	"net"
	"net/http"
)

func runServer(lt net.Listener) {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintf(w, "Hello P2PWN-Go")
	})

	server := http.Server{
		Addr: ":" + Config.Port,
	}

	fmt.Printf("Server is listening on %v\n", server.Addr)
	server.Serve(lt)
}
