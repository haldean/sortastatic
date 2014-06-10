package main

import (
	"log"
	"path"
	"text/template"
	"time"
)

var indexTemplate *template.Template
var pageTemplate *template.Template
var notFoundTemplate *template.Template

func FormatDate(t time.Time) string {
	return t.Format("2006-01")
}

func LoadTemplates() {
	var tpath = func(template string) string {
		return path.Join(*templates, template + ".html")
	}
	var err error

	funcs := make(map[string]interface{})
	funcs["formatDate"] = FormatDate

	indexTemplate = template.Must(
		template.New("index.html").Funcs(funcs).ParseFiles(tpath("index")))
	pageTemplate = template.Must(
		template.New("page.html").Funcs(funcs).ParseFiles(tpath("page")))

	notFoundTemplate, err =
		template.New("404.html").Funcs(funcs).ParseFiles(tpath("404"))
	if err != nil {
		log.Printf("warning: could not load 404 template: %v", err)
	}
}
