package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var clientCh = make(chan *message)

func clientConnect(rawurl string) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// url, urlErr := url.Parse(rawurl)
	url, urlErr := url.Parse("http://localhost:3000/ws")
	if urlErr != nil {
		fmt.Printf("Error Parsing URL: %s\n", urlErr)
		return
	}
	url.Scheme = "ws" // coerce into wss://
	url.Path = "ws"

	c, _, dialErr := websocket.DefaultDialer.Dial(url.String(), nil)
	if dialErr != nil {
		fmt.Printf("dial: %v\n", dialErr)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, rawMessage, mErr := c.ReadMessage()
			if mErr != nil {
				fmt.Printf("read: %v\n", mErr)
				return
			}
			var m = &message{}
			jsonErr := json.Unmarshal(rawMessage, m)
			if jsonErr != nil {
				fmt.Printf("WSS Json Unmarshal Error: %+v\n", jsonErr)
				break
			}

			fmt.Printf("Received message:%+v\n", m)
			m.ActionSquare.set(m.BlockId)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := c.WriteMessage(0, []byte(t.String()))
			if err != nil {
				fmt.Println("write:", err)
				return
			}
		case <-interrupt:
			fmt.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				fmt.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}

}
