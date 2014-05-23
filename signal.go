package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

var pidfile = flag.String(
		"pidfile", "/tmp/sortastatic.pid", "file to write PID to")

func WritePidfile() {
	f, err := os.Create(*pidfile)
	if err != nil {
		log.Printf("warning: can't write pidfile, -reload flag won't work.")
		return
	}
	defer f.Close()
	f.Write([]byte(fmt.Sprintf("%d\n", os.Getpid())))
}

func RegisterSignalHandlers() {
	log.Printf("registering process handlers for pid=%d", os.Getpid())
	WritePidfile()
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

func SendReload() {
	f, err := os.Open(*pidfile)
	if err != nil {
		log.Printf("couldn't open pidfile: %v", err)
		return
	}
	buf := make([]byte, 32)
	n, err := f.Read(buf)
	if err != nil {
		log.Printf("couldn't read pidfile: %v", err)
		return
	}
	pidstr := strings.Trim(string(buf[:n]), " \n")
	pid, err := strconv.ParseInt(pidstr, 10, 32)
	if err != nil {
		log.Printf("couldn't decode pidfile: %v", err)
		return
	}
	err = syscall.Kill(int(pid), syscall.SIGUSR1)
	if err != nil {
		log.Printf("couldn't send reload signal: %v", err)
	}
}
