package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

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
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(client); err != nil {
		fmt.Printf("Error decoding response: %s", err)
	}

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
			defer resp.Body.Close()

			decoder := json.NewDecoder(resp.Body)
			if err := decoder.Decode(world); err != nil {
				fmt.Printf("Error decoding response: %s", err)
			}
		}

	}()
}

func clientAction(ca *ClientAction) {
	go func() {
		ca.ClientID = client.ID
		jsonCa, _ := json.Marshal(ca)

		_, reqErr := http.Post(clientCreateUrl("/srv/action"), "application/json", bytes.NewBuffer(jsonCa))
		if reqErr != nil {
			fmt.Printf("Error in http request, client action: %s\n", reqErr)
			return
		}
	}()
}
