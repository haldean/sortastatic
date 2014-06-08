package main

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type CompressWriter struct {
	inner http.ResponseWriter
	writer *gzip.Writer
}

func NewWriter(inner http.ResponseWriter) CompressWriter {
	return CompressWriter{inner: inner, writer: gzip.NewWriter(inner)}
}

func (c *CompressWriter) Header() http.Header {
	return c.inner.Header()
}

func (c *CompressWriter) WriteHeader(status int) {
	c.inner.WriteHeader(status)
}

func (c *CompressWriter) Write(data []byte) (int, error) {
	c.Header().Add("Content-Encoding", "gzip")
	c.Header().Add("Vary", "Accept-Encoding")
	if c.Header().Get("Content-Type") == "" {
		// otherwise it'll sniff it as gzip, which we Do Not Want
		c.Header().Set("Content-Type", http.DetectContentType(data))
	}
	return c.writer.Write(data)
}

func (c *CompressWriter) Close() {
	c.writer.Close()
}

type CompressHandler struct {
	inner http.Handler
}

func canUseGzip(r http.Request) bool {
	encodings, ok := r.Header["Accept-Encoding"]
	if !ok {
		return false;
	}
	for _, encs := range encodings {
		es := strings.Split(encs, ",")
		for _, e := range es {
			if strings.ToLower(e) == "gzip" {
				return true;
			}
		}
	}
	return false;
}

func Compress(handler http.Handler) http.Handler {
	return CompressHandler{inner: handler}
}

func (c CompressHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if (canUseGzip(*r)) {
		cw := NewWriter(w)
		defer cw.Close()
		c.inner.ServeHTTP(&cw, r)
	} else {
		c.inner.ServeHTTP(w, r)
	}
}
