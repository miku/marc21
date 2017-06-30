// MARC21 to XML Converter written in Go.
//
// Copyright (C) 2011 William Waites
// Copyright (C) 2017 Martin Czygan, <martin.czygan@uni-leipzig.de>
//
// This program is free software: you can redistribute it and/or
// modify it under the terms of the GNU General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// General Public License for more details.

// You should have received a copy of the GNU General Public License
// and the GNU General Public License along with this program (the
// named GPL3).  If not, see <http://www.gnu.org/licenses/>.
package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/miku/marc21"
)

var declaration = `<?xml version="1.0" encoding="utf-8" ?>`

// stickyErrWriter keeps an error around, so you can *occasionally* check if an error occured.
type stickyErrWriter struct {
	w   io.Writer
	err *error
}

func (sew stickyErrWriter) Write(p []byte) (n int, err error) {
	if *sew.err != nil {
		return 0, *sew.err
	}
	n, err = sew.w.Write(p)
	*sew.err = err
	return
}

func main() {
	flag.Parse()

	var rc = ioutil.NopCloser(os.Stdin)
	var writer io.Writer = os.Stdout
	var err error

	if flag.NArg() > 0 {
		if rc, err = os.Open(flag.Arg(0)); err != nil {
			log.Fatal(err)
		}
		defer rc.Close()
	}

	if flag.NArg() > 1 {
		if writer, err = os.Open(flag.Arg(1)); err != nil {
			log.Fatal(err)
		}
	}

	w := &stickyErrWriter{writer, &err}

	var once sync.Once

	for {
		record, err := marc21.ReadRecord(rc)
		if err == io.EOF {
			break
		}
		once.Do(func() {
			io.WriteString(w, declaration)
			io.WriteString(w, `<collection xmlns="http://www.loc.gov/MARC21/slim">`)
		})
		if err != nil {
			log.Fatal(err)
		}
		record.WriteTo(w)
	}
	io.WriteString(w, "</collection>\n")
	if *w.err != nil {
		log.Fatal(err)
	}
}
