package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/kr/pretty"

	"github.com/twpayne/go-gpx"
)

func dumpFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return dump(f)
}

func dump(r io.Reader) error {
	g, err := gpx.Read(r)
	if err != nil {
		return err
	}
	pretty.Println(g)
	return nil
}

func run() error {
	flag.Parse()
	if flag.NArg() == 0 {
		return dump(os.Stdin)
	}
	for _, arg := range flag.Args() {
		if err := dumpFile(arg); err != nil {
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
