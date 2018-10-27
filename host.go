package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"

	"github.com/localtunnel/go-localtunnel"
)

// CHANGE ME
const appName = "p2pwn-go-maze"
const appRelease = "DEVELOPMENT"

// App  Config
var Config = &appConfig{}

type appConfig struct {
	AppName string `json:"appName"` // Name of this app
	Release string `json:"release"` // "PRODUCTION", "DEVELOPMENT"
	Port    string // Server Listening Port
	P2pwn   string // P2PWN Service Address
}

// P2PWN Service Config
var P2pwn = &p2pwnConfig{}

type p2pwnConfig struct { // all values will be provided by P2PWN
	ID          string `json:"id"`           // public id assigned by P2PWN service
	AccessToken string `json:"access_token"` // private access token needed to perform actions on this host
	DisplayName string `json:"display_name"` // name supplied by appName for grouping rooms in P2PWN
	EntryURL    string `json:"entry_url"`    // url used as the entrypoint for your app, supplied by localtunnel
}

func setConfig(configPtr *string, flagName string, defaultVal string, help string) {
	flag.StringVar(configPtr, flagName, defaultVal, help)

	if val, ok := os.LookupEnv(flagName); ok {
		*configPtr = val
	}
}

func structToMap(i interface{}) (values url.Values) {
	values = url.Values{}
	iVal := reflect.ValueOf(i).Elem()
	typ := iVal.Type()
	for i := 0; i < iVal.NumField(); i++ {
		values.Set(typ.Field(i).Name, fmt.Sprint(iVal.Field(i)))
	}
	return
}

func host() {

	setConfig(&Config.AppName, "name", appName, "Name of this app")
	setConfig(&Config.Port, "port", "3000", "Port for server to listen on")
	setConfig(&Config.P2pwn, "p2pwn", "https://p2pwithme.2018.nodeknockout.com", "P2PWN Service Address")

	flag.Parse()

	port, portErr := strconv.Atoi(Config.Port)
	if portErr != nil {
		fmt.Printf("Invalid Port config: %s -> %v \n", Config.Port, port)
		os.Exit(1)
		return
	}

	lt, ltErr := localtunnel.Listen(localtunnel.Options{
		Subdomain: Config.AppName,
	})
	if ltErr != nil {
		fmt.Printf("Error creating localtunnel: %v\n", ltErr)
		os.Exit(1)
		return
	}

	p2pwnRes, p2pwnErr := http.PostForm(Config.P2pwn, structToMap(Config))
	if p2pwnErr != nil {
		fmt.Printf("Error Connecting to P2PWN Service: %v\n", p2pwnErr)
		os.Exit(1)
		return
	}

	defer p2pwnRes.Body.Close()
	if err := json.NewDecoder(p2pwnRes.Body).Decode(P2pwn); err != nil {
		fmt.Println("Unmarshal P2PWN Response Error:", err)
		os.Exit(1)
		return
	}

	fmt.Printf("P2PWN is Ready: %+v\n", P2pwn)

	// Setup your handlers
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
