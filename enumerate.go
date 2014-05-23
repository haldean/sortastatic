package main

import (
	"fmt"
	"log"
	"path/filepath"
	"sort"
)

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
