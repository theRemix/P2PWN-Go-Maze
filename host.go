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

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/localtunnel/go-localtunnel"
	"golang.org/x/image/colornames"
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

func runHost() {
	const font = "fonts/zorque.ttf"

	cfg := pixelgl.WindowConfig{
		Title:  "Host Game",
		Bounds: pixel.R(0, 0, 1024, 768),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)

	fontFace, fontFaceErr := loadTTF(font, 80)
	if fontFaceErr != nil {
		panic(fontFaceErr)
	}

	win.Clear(colornames.Firebrick)
	atlas := text.NewAtlas(fontFace, text.ASCII)

	titleTxt := text.New(pixel.V(350, 100), atlas)
	titleTxt.Color = colornames.Lightgrey
	titleTxt.WriteString("Host Game")
	titleTxt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(pixel.V(600, -100))))

	statusTxt := text.New(pixel.V(350, 100), atlas)
	statusTxt.Color = colornames.Darkkhaki
	statusTxt.WriteString("Creating Server")
	statusTxt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(pixel.V(730, 200))))

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

	fmt.Printf("Connected to LT: %v\n", lt.URL())

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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintf(w, "Hello P2PWN-Go")
	})

	server := http.Server{
		Addr: ":" + Config.Port,
	}

	fmt.Printf("Server is listening on %v\n", server.Addr)
	go server.Serve(lt)

	statusTxt.Clear()
	statusTxt.Color = colornames.Darkcyan
	statusTxt.WriteString("Connected!")
	statusTxt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(pixel.V(630, 300))))

	for !win.Closed() {
		if win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ) {
			server.Shutdown(nil)
			go func() { exitCh <- true }()
			return
		}

		win.Update()
	}

}
