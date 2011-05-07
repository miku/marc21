package main

import (
	"bitbucket.org/ww/marc21"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var gzoutput bool

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n%s\n\n", "Go MARC21 - XML Converter")
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] infile.mrc [outfile.xml]\n\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}
	flag.BoolVar(&gzoutput, "z", false, "Gzip Output")
}

func main() {
	flag.Parse()

	nargs := flag.NArg()
	if flag.NArg() == 0 || flag.NArg() > 2 {
		flag.Usage()
		os.Exit(255)
	}

	var infile io.Reader
	var outfile io.Writer

	infile, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	if nargs == 2 {
		outfile, err = os.Open(flag.Arg(1))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		outfile = os.Stdout
	}

	if gzoutput {
		outfile, err = gzip.NewWriter(outfile)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = marc2xml(infile, outfile)
	if err != nil {
		log.Fatal(err)
	}

	closer, ok := infile.(io.Closer)
	if ok {
		closer.Close()
	}

	closer, ok = outfile.(io.Closer)
	if ok {
		closer.Close()
	}
}

func marc2xml(reader io.Reader, writer io.Writer) (err os.Error) {
	records := make(chan *marc21.Record)

	go func() {
		defer close(records)
		for {
			record, err := marc21.ReadRecord(reader)
			if err == os.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			records <- record
		}
	}()

	_, err = writer.Write([]byte(`<?xml version="1.0" encoding="utf-8" ?>
<collection xmlns="http://www.loc.gov/MARC21/slim">\n`))
	if err != nil {
		return
	}

	for record := range records {
		err = record.XML(writer)
		if err != nil {
			return
		}
	}

	_, err = writer.Write([]byte("</collection>\n"))
	return
}
