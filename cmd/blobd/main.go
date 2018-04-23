package main

import (
	"flag"

	"github.com/reflect/xflag"

	"github.com/bradhe/blobd/logs"
	"github.com/bradhe/blobd/server"
)

func main() {
	var (
		listenAddr = flag.String("listen-addr", "localhost:8765", "address to bind service to")
		storageURL = xflag.URL("storage-url", "", "URL that storage will be stored in")
		debug      = flag.Bool("debug", false, "put the server in debug mode")
	)

	// TODO: Validate flags.
	flag.Parse()

	opts := server.ServerOptions{
		StorageURL: *storageURL,
	}

	if *debug {
		logs.EnableDebug()
	} else {
		logs.DisableDebug()
	}

	logs.WithPackage("main").Printf("starting blobd v%s", Version)
	logs.WithPackage("main").Printf("listening on addr %s", *listenAddr)

	s := server.New(opts)
	s.ListenAndServe(*listenAddr)
}
