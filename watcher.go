package main

import (
	"fmt"
	"log"
	"path/filepath"
	"golang.org/x/exp/inotify"
)

func Watch() {
	watcher, err := inotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	paths, err := filepath.Glob(fmt.Sprintf("%s/*/*", *index))
	if err != nil {
		log.Fatal(err)
	}

	for _, path := range paths {
		fmt.Printf("watching %s\n", path)
		err = watcher.AddWatch(path, inotify.IN_MODIFY | inotify.IN_DELETE_SELF)
		if err != nil {
			log.Fatal(err)
		}
	}

	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if ev == nil {
					continue
				}

				log.Println("file was updated, reloading:", ev)
				LoadTemplates()
				BuildCache()

				if ev.Mask & inotify.IN_DELETE_SELF != 0 {
					log.Println("file was deleted, restartig:", ev)
					watcher.Close()
					Watch()
				}
			case <-watcher.Error:
				break
			}
		}
	}()
}
