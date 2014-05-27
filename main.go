package main

import (
	"flag"
	"log"
	"net/http"
	"path"
	"text/template"
)

var index = flag.String("index", "", "location of markdown files")
var templates = flag.String("templates", "", "location of page template")
var static = flag.String("static", "", "location of files shared between pages")
var reload = flag.Bool("reload", false, "reload a running instance of sortastatic")
var bindAddr = flag.String("bind", "127.0.0.1:8080", "address/port to bind to")

var pages map[string]*Page
var sorted []*Page
var indexTemplate *template.Template
var pageTemplate *template.Template
var notFoundTemplate *template.Template

func LoadTemplates() {
	var tpath = func(template string) string {
		return path.Join(*templates, template + ".html")
	}

	var err error
	indexTemplate, err = template.ParseFiles(tpath("index"))
	if err != nil {
		log.Fatalf("could not load index template: %v", err)
	}
	pageTemplate, err = template.ParseFiles(tpath("page"))
	if err != nil {
		log.Fatalf("could not load page template: %v", err)
	}
	notFoundTemplate, err = template.ParseFiles(tpath("404"))
	if err != nil {
		log.Printf("warning: could not load 404 template: %v", err)
	}
}

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

	log.Printf("binding to %s", *bindAddr)
	log.Fatal(http.ListenAndServe(*bindAddr, NewMux(*static)))
}
