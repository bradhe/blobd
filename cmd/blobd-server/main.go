package main

import (
	"flag"
	"github.com/bradhe/blobd/server"
	"github.com/reflect/xflag"
	"log"
)

func main() {
	var (
		listenAddr = flag.String("listen-addr", "localhost:8765", "address to bind service to")
		storageURL = xflag.URL("storage-url", "", "URL that storage will be stored in")
	)

	// TODO: Validate flags.
	flag.Parse()

	opts := server.ServerOptions{
		StorageURL: *storageURL,
	}

	log.Printf("starting blobd v%s", Version)
	log.Printf("listening on addr %s", *listenAddr)

	s := server.New(opts)
	s.ListenAndServe(*listenAddr)
}
