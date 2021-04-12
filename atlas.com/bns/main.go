package main

import (
	"atlas-bns/name"
	"atlas-bns/rest"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	l := log.New(os.Stdout, "bns ", log.LstdFlags|log.Lmicroseconds)

	name.InitCache(l)

	go rest.GetServer(l).Run()

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	l.Println("[INFO] Shutting down via signal:", sig)
}
