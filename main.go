package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
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

func Write404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	if notFoundTemplate != nil {
		notFoundTemplate.Execute(w, r)
	} else {
		w.Write([]byte("four oh four, yo. four. zero. four."))
	}
}

func Write500(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("five hundred! shit! five! hundred!"))
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	indexTemplate.Execute(w, sorted)
}

// PageHandler is a Handler instead of a HandlerFunc so we can wrap it in a
// StripPrefix.
type PageHandler struct{}

func (_ PageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(path) < 1 {
		Write404(w, r)
		return
	}
	p, ok := pages[path[0]]
	if !ok {
		Write404(w, r)
	} else {
		if r.URL.RawQuery == "raw" {
			f, err := os.Open(p.Body)
			if err != nil {
				Write500(w, r)
				return
			}
			_, err = io.Copy(w, f)
			if err != nil {
				Write500(w, r)
			}
			return
		}

		if len(path) == 1 {
			p.Load()
			pageTemplate.Execute(w, p)
		}
	}
}

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

	http.HandleFunc("/", IndexHandler)
	http.Handle("/x/", http.StripPrefix("/x/", PageHandler{}))
	if *commondir != "" {
		log.Printf("using %v as the common directory", *commondir)
		http.Handle("/c/", http.StripPrefix("/c/",
			http.FileServer(http.Dir(*commondir))))
	}
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
