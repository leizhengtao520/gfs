package main

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 1 {
		usage()
		os.Exit(2)
	}
	mountpoint := flag.Arg(0)
	c, err := fuse.Mount(
		mountpoint,
		fuse.FSName("leiFs"),
		fuse.Subtype("leiFs"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	err = fs.Serve(c, FS{})
	if err != nil {
		log.Fatal(err)
	}
}
func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "   %s MOUNTPOINT", os.Args[0])
	flag.PrintDefaults()
}
