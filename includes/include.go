package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// Reads all .txt files in the current folder
// and encodes them as strings literals in textfiles.go
func main() {
	fs, _ := ioutil.ReadDir("./includes")
	out, _ := os.Create("inlined.go")
	out.Write([]byte("package main \n\nconst (\n"))
	for _, f := range fs {
		if strings.HasSuffix(f.Name(), ".go") {
			break
		}
		if strings.HasSuffix(f.Name(), ".html") {
			out.Write([]byte(strings.TrimSuffix(f.Name(), ".html") + " = `"))
			f, _ := os.Open("./includes/" + f.Name())
			io.Copy(out, f)
			out.Write([]byte("`\n"))
		}
		if strings.HasSuffix(f.Name(), ".ttf") {

			out.Write([]byte(strings.TrimSuffix(f.Name(), ".ttf") + " = `"))
			f, _ := os.Open("./includes/" + f.Name())
			bytes, _ := ioutil.ReadAll(f)
			fmt.Fprintf(out, "%X", bytes)
			// out.Write(bytes)
			// io.Copy(out, f)
			out.Write([]byte("`\n"))
		}
	}
	out.Write([]byte(")\n"))
}
