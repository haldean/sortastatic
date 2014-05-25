package main

import (
	"flag"
	"log"
	"net/http"
	"text/template"
)

var index = flag.String("index", "", "location of markdown files")
var itemplate = flag.String("index-template", "", "location of index template")
var nftemplate = flag.String("404-template", "", "location of 404 page template")
var ptemplate = flag.String("page-template", "", "location of page template")
var commondir = flag.String("common-dir", "", "location of files shared between pages")
var reload = flag.Bool("reload", false, "reload a running instance of sortastatic")

var pages map[string]*Page
var sorted []*Page
var indexTemplate *template.Template
var pageTemplate *template.Template
var notFoundTemplate *template.Template

func LoadTemplates() {
	var err error
	indexTemplate, err = template.ParseFiles(*itemplate)
	if err != nil {
		log.Fatalf("could not load index template: %v", err)
	}
	notFoundTemplate, err = template.ParseFiles(*nftemplate)
	if err != nil {
		log.Printf("warning: could not load 404 template: %v", err)
	}
	pageTemplate, err = template.ParseFiles(*ptemplate)
	if err != nil {
		log.Fatalf("could not load page template: %v", err)
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
	if *itemplate == "" {
		log.Fatal("must specify index template location")
	}
	if *ptemplate == "" {
		log.Fatal("must specify page template location")
	}

	RegisterSignalHandlers()
	LoadTemplates()
	BuildCache()

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", NewMux(*commondir)))
}
