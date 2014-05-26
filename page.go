package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/russross/blackfriday"
)

type Page struct {
	Name     string
	Path     string
	Title    string
	Body     string
	Url      string
	Rendered string
	UseMarkdown bool
	Public   bool
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

	index := fmt.Sprintf("%s/index.html", p.Path)
	if FileExists(index) {
		p.UseMarkdown = false
		p.Body = index
	} else {
		mds, err := filepath.Glob(fmt.Sprintf("%s/*.md", p.Path))
		if err != nil {
			return p, err
		}
		if len(mds) == 1 {
			p.Body = mds[0]
			p.UseMarkdown = true
		} else if len(mds) == 0 {
			return p, errors.New(
				fmt.Sprintf("no markdown files or index for page %v", p.Name))
		} else if len(mds) > 1 {
			log.Printf("warning: more than one markdown file found for page %v, "+
				"using %v", p.Name, mds[0])
			p.Body = mds[0]
			p.UseMarkdown = true
		}
	}

	err = p.LoadTitle()
	if err != nil {
		return p, err
	}

	p.Public = !FileExists(fmt.Sprintf("%s/draft", p.Path))

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
		p.Rendered = string(blackfriday.MarkdownBasic(data))
	} else {
		p.Rendered = string(data)
	}
	return nil
}

type ByTitle []*Page

func (a ByTitle) Len() int           { return len(a) }
func (a ByTitle) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTitle) Less(i, j int) bool { return a[i].Title < a[j].Title }
