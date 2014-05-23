package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/russross/blackfriday"
)

var index = flag.String("index", "", "location of markdown files")
var itemplate = flag.String("index-template", "", "location of index template")
var nftemplate = flag.String("404-template", "", "location of 404 page template")
var ptemplate = flag.String("page-template", "", "location of page template")
var commondir = flag.String("common-dir", "", "location of files shared between pages")

type Page struct {
	Name  string
	Path  string
	Title string
	Body  string
	Url   string
	Rendered string
}

func NewPage(path string) (Page, error) {
	var p Page
	var err error
	p.Path, err = filepath.Abs(path)
	if err != nil {
		return p, err
	}
	p.Name = filepath.Base(path)
	p.Url = fmt.Sprintf("x/%s", p.Name)

	mds, err := filepath.Glob(fmt.Sprintf("%s/*.md", p.Path))
	if err != nil {
		return p, err
	}
	if len(mds) == 0 {
		return p, errors.New(
			fmt.Sprintf("no markdown files for page %v", p.Name))
	} else if len(mds) > 1 {
		log.Printf("warning: more than one markdown file found for page %v, "+
			"using %v", p.Name, mds[0])
	}
	p.Body = mds[0]

	f, err := os.Open(p.Body)
	if err != nil {
		return p, err
	}

	var i int
	buf := make([]byte, 256)
	n, err := f.Read(buf)
	if err != nil {
		return p, err
	}
	for i = 0; i < n; i++ {
		if buf[i] == '\n' {
			break
		}
	}
	p.Title = string(buf[:i])

	return p, nil
}

func (p *Page) Load() error {
	if len(p.Rendered) != 0 {
		return nil;
	}
	log.Printf("loading content for %v", p.Name)

	data, err := ioutil.ReadFile(p.Body)
	if err != nil {
		return err
	}
	p.Rendered = string(blackfriday.MarkdownCommon(data))
	return nil
}

type ByTitle []*Page

func (a ByTitle) Len() int           { return len(a) }
func (a ByTitle) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTitle) Less(i, j int) bool { return a[i].Title < a[j].Title }

var pages map[string]*Page
var sorted []*Page
var indexTemplate *template.Template
var pageTemplate *template.Template
var notFoundTemplate *template.Template

func BuildCache() {
	pages = make(map[string]*Page)
	sorted = make([]*Page, 0)

	paths, err := filepath.Glob(fmt.Sprintf("%s/*", *index))
	if err != nil {
		log.Printf("error when finding files: %v", err)
	}
	for i := 0; i < len(paths); i++ {
		p, err := NewPage(paths[i])
		if err != nil {
			log.Printf("page could not be loaded: %v", err)
		} else {
			pages[p.Name] = &p
			sorted = append(sorted, &p)
			log.Printf("loaded: %v", p.Title)
		}
	}
	sort.Sort(ByTitle(sorted))
	log.Printf("loaded %d pages", len(pages))
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	indexTemplate.Execute(w, sorted)
}

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
	if *index == "" {
		log.Fatal("must specify index location")
	}
	if *itemplate == "" {
		log.Fatal("must specify index template location")
	}
	if *ptemplate == "" {
		log.Fatal("must specify page template location")
	}

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
