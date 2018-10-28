package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// Encodes files as strings literals in textfiles.go
func main() {
	fs, _ := ioutil.ReadDir("./includes")
	out, _ := os.Create("inlined.go")
	out.Write([]byte("package main \n\nconst (\n"))

	out.Write([]byte("index = `"))
	f, _ := os.Open("./static/index.html")
	io.Copy(out, f)
	out.Write([]byte("`\n"))

	for _, f := range fs {
		if strings.HasSuffix(f.Name(), ".go") {
			break
		}
		if strings.HasSuffix(f.Name(), ".ttf") {

			out.Write([]byte(strings.TrimSuffix(f.Name(), ".ttf") + " = `"))
			f, _ := os.Open("./includes/" + f.Name())
			bytes, _ := ioutil.ReadAll(f)
			fmt.Fprintf(out, "%X", bytes)
			out.Write([]byte("`\n"))
		}
	}
	out.Write([]byte(")\n"))
}
