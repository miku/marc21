package main

import (
	"io"
	"log"
	"os"

	"github.com/miku/marc21"
)

func main() {
	io.WriteString(os.Stdout, "<collection>")
	for {
		r, err := marc21.ReadRecord(os.Stdin)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if _, err := r.WriteTo(os.Stdout); err != nil {
			log.Fatal(err)
		}
	}
	io.WriteString(os.Stdout, "</collection>")
}
