package main

import (
	"mime"
	"path/filepath"
	"strings"
)

func MimeType(path string) string {
	return mime.TypeByExtension(filepath.Ext(path))
}

func StripHtmlComments(line string) string {
	if !strings.HasPrefix(line, "<!--") || !strings.HasSuffix(line, "-->") {
		return line
	}
	return strings.TrimSpace(line[4 : len(line)-3])
}
