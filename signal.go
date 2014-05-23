package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func RegisterSignalHandlers() {
	log.Printf("registering process handlers for pid=%d", os.Getpid())
	c := make(chan os.Signal, 1)
	go func() {
		for _ = range c {
			log.Printf("got signal to reload, doing it now")
			LoadTemplates()
			BuildCache()
		}
	}()
	signal.Notify(c, syscall.SIGUSR1)
}
