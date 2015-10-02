package main

import (
	"flag"
	"log"
	"net/http"
)

var index = flag.String("index", "", "location of markdown files")
var templates = flag.String("templates", "", "location of page template")
var static = flag.String("static", "", "location of files shared between pages")
var reload = flag.Bool("reload", false, "reload a running instance of sortastatic")
var watch = flag.Bool("watch", true, "watch for file changes and reload automatically")
var bindAddr = flag.String("bind", "127.0.0.1:8080", "address/port to bind to")

var pages map[string]*Page
var sorted []*Page

func main() {
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

	log.Printf("binding to %s", *bindAddr)
	log.Fatal(http.ListenAndServe(*bindAddr, Compress(NewMux(*static))))
}
