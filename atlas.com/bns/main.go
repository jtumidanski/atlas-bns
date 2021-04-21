package main

import (
	"atlas-bns/logger"
	"atlas-bns/name"
	"atlas-bns/rest"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	l := logger.CreateLogger()

	name.InitCache(l)

	go rest.GetServer(l).Run()

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	l.Infoln("Shutting down via signal:", sig)
}
