package main

import (
	"atlas-cts/configuration"
	"atlas-cts/logger"
	"atlas-cts/rest"
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	l := logger.CreateLogger()
	l.Infoln("Starting main service.")

	config, err := configuration.NewConfigurator(l).GetConfiguration()
	if err != nil {
		l.WithError(err).Errorf("Unable to load service configuration.")
		return
	}


	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	rest.CreateRestService(l, config, ctx, wg)

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	l.Infof("Initiating shutdown with signal %s.", sig)
	cancel()
	wg.Wait()
	l.Infoln("Service shutdown.")
}
