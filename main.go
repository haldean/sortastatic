package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"github.com/lemenkov/systemd.go"
)

var index = flag.String("index", "", "location of markdown files")
var templates = flag.String("templates", "", "location of page template")
var static = flag.String("static", "", "location of files shared between pages")
var reload = flag.Bool("reload", false, "reload a running instance of sortastatic")
var watch = flag.Bool("watch", false, "watch for file changes and reload automatically")

var pages map[string]*Page
var sorted []*Page

func main() {
	var listener net.Listener
	var err error

	flag.Parse()
	if *reload {
		SendReload()
		return
	}

	if *index == "" {
		log.Fatal("must specify index location")
	}
	if *templates == "" {
		log.Fatal("must specify index template location")
	}

	RegisterSignalHandlers()
	LoadTemplates()
	BuildCache()

	if *watch {
		Watch()
	}

	sockets := systemd.ListenFds()
	if sockets == nil {
		listener, err = net.Listen("unix", "/tmp/sortastatic.socket")
	} else {
		listener, err = net.FileListener(sockets[0])
	}
	if err != nil {
		panic(err);
	}
	log.Fatal(http.Serve(listener, Compress(NewMux(*static))))
}
