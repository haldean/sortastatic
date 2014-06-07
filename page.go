package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/haldean/gommd"
)

const mmdflags = gommd.EXT_PROCESS_HTML | gommd.EXT_SMART

type Page struct {
	Name        string
	Path        string
	Title       string
	Body        string
	Url         string
	Rendered    string
	Stylesheets []string
	UseMarkdown bool
	Public      bool
}

func NewPage(path string) (Page, error) {
	var p Page
	var err error
	p.Path, err = filepath.Abs(path)
	if err != nil {
		return p, err
	}
	p.Name = filepath.Base(path)
	p.Url = fmt.Sprintf("%s/", p.Name)
	p.Public = !FileExists(fmt.Sprintf("%s/draft", p.Path))

	log.Printf("loading %s", p.Name)

	err = p.FindBody()
	if err != nil {
		return p, err
	}

	err = p.LoadTitle()
	if err != nil {
		return p, err
	}

	return p, nil
}

func (p *Page) LoadTitle() error {
	f, err := os.Open(p.Body)
	if err != nil {
		return err
	}
	defer f.Close()

	var i int
	buf := make([]byte, 256)
	n, err := f.Read(buf)
	if err != nil {
		return err
	}
	for i = 0; i < n; i++ {
		if buf[i] == '\n' {
			break
		}
	}
	p.Title = StripHtmlComments(string(buf[:i]))
	return nil
}

func (p *Page) FindBody() error {
	index := fmt.Sprintf("%s/index.html", p.Path)
	if FileExists(index) {
		log.Printf("  using index file")
		p.UseMarkdown = false
		p.Body = index
		return nil
	}

	mds, err := filepath.Glob(fmt.Sprintf("%s/*.md", p.Path))
	if err != nil {
		return err
	}
	if len(mds) == 1 {
		p.Body = mds[0]
		p.UseMarkdown = true
	} else if len(mds) == 0 {
		return errors.New(
			fmt.Sprintf("no markdown files or index for page %v", p.Name))
	} else if len(mds) > 1 {
		log.Printf("warning: more than one markdown file found for page "+
			"%v, using %v", p.Name, mds[0])
		p.Body = mds[0]
		p.UseMarkdown = true
	}

	p.Stylesheets = make([]string, 0)
	ss, err := filepath.Glob(fmt.Sprintf("%s/*.css", p.Path))
	if err != nil {
		log.Printf("  stylesheets couldn't be loaded: %v", err)
	} else {
		log.Printf("  found %d stylesheets", len(ss))
		for _, s := range ss {
			p.Stylesheets = append(p.Stylesheets, filepath.Base(s))
		}
	}

	return nil
}

func (p *Page) Load() error {
	if len(p.Rendered) != 0 {
		return nil
	}
	log.Printf("loading content for %v", p.Name)

	data, err := ioutil.ReadFile(p.Body)
	if err != nil {
		return err
	}
	if p.UseMarkdown {
		p.Rendered = gommd.MarkdownToString(
			string(data), mmdflags, gommd.FORMAT_HTML)
	} else {
		p.Rendered = string(data)
	}
	return nil
}

type ByTitle []*Page

func (a ByTitle) Len() int      { return len(a) }
func (a ByTitle) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByTitle) Less(i, j int) bool {
	return strings.ToLower(a[i].Title) < strings.ToLower(a[j].Title)
}
