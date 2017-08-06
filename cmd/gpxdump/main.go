package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/twpayne/go-gpx"
)

func dumpFile(w io.Writer, filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return dump(w, f)
}

func dump(w io.Writer, r io.Reader) error {
	g, err := gpx.Read(r)
	if err != nil {
		return err
	}
	spew.Fdump(w, g)
	return nil
}

func run() error {
	flag.Parse()
	if flag.NArg() == 0 {
		return dump(os.Stdout, os.Stdin)
	}
	for _, arg := range flag.Args() {
		if err := dumpFile(os.Stdout, arg); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
