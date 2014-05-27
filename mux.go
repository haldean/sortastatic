package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

type Mux struct {
	commonHandler http.Handler
}

func NewMux(commonPath string) Mux {
	var m Mux
	if commonPath != "" {
		m.commonHandler = http.StripPrefix(
			"/c/", http.FileServer(http.Dir(commonPath)))
	}
	return m
}

func (m Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s : %s", r.Method, r.URL, r.Header.Get("User-Agent"))
	if r.URL.Path == "/" {
		IndexHandler(w, r)
		return
	}
	if PageHandler(w, r) {
		return
	}
	if strings.HasPrefix(r.URL.Path, "/c/") {
		m.commonHandler.ServeHTTP(w, r)
		return
	}
	Write404(w, r)
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

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	indexTemplate.Execute(w, sorted)
}

// returns true if this handler should be handling the request, even if the
// request is invalid (i.e., returns 404/500)
func PageHandler(w http.ResponseWriter, r *http.Request) bool {
	urlpath := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(urlpath) < 1 {
		return false
	}
	p, ok := pages[urlpath[0]]
	if !ok {
		return false
	}

	if len(urlpath) == 1 && !strings.HasSuffix(r.URL.Path, "/") {
		http.Redirect(w, r, r.URL.Path+"/", http.StatusMovedPermanently)
		return true
	}

	if len(urlpath) == 1 {
		if r.URL.RawQuery == "raw" || !p.UseMarkdown {
			ServeFile(p.Body, w, r)
		} else {
			p.Load()
			pageTemplate.Execute(w, p)
		}
		return true
	}

	// check for static page file
	urlpath = urlpath[1:]
	fpath := path.Join(p.Path, path.Join(urlpath...))
	if !FileExists(fpath) {
		Write404(w, r)
		return true
	}
	ServeFile(fpath, w, r)

	return true
}

func ServeFile(fpath string, w http.ResponseWriter, r *http.Request) {
	mimetype := MimeType(fpath)
	if mimetype != "" {
		w.Header().Add("Content-type", mimetype)
	}

	f, err := os.Open(fpath)
	if err != nil {
		Write500(w, r)
		return
	}
	_, err = io.Copy(w, f)
	if err != nil {
		Write500(w, r)
	}
}
