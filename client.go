package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type ClientAction struct {
	ClientID     int
	Op           OpCode
	ActionSquare actionSquare
	BlockId      int
}

type Client struct {
	ID      int
	Updated time.Time
}

func (c *Client) NeedsUpdate() bool {
	return !c.Updated.Equal(worldState.LastUpdated)
}

func (c *Client) SendAction(ca *ClientAction) {
	if c.ID == 0 { // client is host
		return
	}
	ca.ClientID = c.ID
	fmt.Printf("Client Action: %+v\n", ca)
	go func() {
		jsonCa, _ := json.Marshal(ca)

		resp, reqErr := http.Post(clientCreateUrl("/srv/action"), "application/json", bytes.NewBuffer(jsonCa))
		if reqErr != nil {
			fmt.Printf("Error in http request, client action: %s\n", reqErr)
			return
		}
		if resp.StatusCode != 200 {
			go func() { stateCh <- Menu }()
			return
		}
	}()
}

type UpdateResponse struct {
	State      OpCode
	Updated    time.Time
	WorldTiles *WorldTiles
}

var (
	client  = &Client{}
	hostURL string
)

func clientCreateUrl(path string) string {
	url, urlErr := url.Parse(hostURL)
	if urlErr != nil {
		fmt.Printf("Error Parsing URL: %s\n", urlErr)
		return ""
	}
	url.Path = path

	return url.String()
}

func clientConnect(rawurl string) {
	hostURL = rawurl

	resp, reqErr := http.Post(clientCreateUrl("/srv/connect"), "application/json", nil)
	if reqErr != nil {
		fmt.Printf("Error in http request, client connect: %s\n", reqErr)
		return
	}
	if resp.StatusCode != 200 {
		go func() { stateCh <- Menu }()
		return
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(client); err != nil {
		fmt.Printf("Error decoding response: %s\n", err)
	}

	fmt.Printf("Client Connected to Host: %s\n", hostURL)

	clientUpdates()
}

func clientUpdates() {
	ping := time.NewTicker(time.Millisecond * 500)
	go func() {
		for _ = range ping.C {

			jsonClient, _ := json.Marshal(client)

			resp, reqErr := http.Post(clientCreateUrl("/srv/update"), "application/json", bytes.NewBuffer(jsonClient))
			if reqErr != nil {
				fmt.Printf("Error in http request, client update: %s\n", reqErr)
				return
			}
			if resp.StatusCode != 200 {
				go func() { stateCh <- Menu }()
				return
			}
			defer resp.Body.Close()

			updateResponse := &UpdateResponse{}

			decoder := json.NewDecoder(resp.Body)
			if err := decoder.Decode(updateResponse); err != nil {
				fmt.Printf("Error decoding response: %s\n", err)
			}

			if updateResponse.State == StateNew {
				client.Updated = updateResponse.Updated
				worldState.Tiles = updateResponse.WorldTiles
			}
		}

	}()
}
