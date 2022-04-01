package main

import (
	"atlas-cts/configuration"
	"atlas-cts/kafka"
	"atlas-cts/logger"
	"atlas-cts/rest"
	"atlas-cts/tasks"
	"atlas-cts/tracing"
	"atlas-cts/transport"
	"context"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const serviceName = "atlas-cts"
const consumerGroupId = "Character Transport Service"

func main() {
	l := logger.CreateLogger(serviceName)
	l.Infoln("Starting main service.")

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	tc, err := tracing.InitTracer(l)(serviceName)
	if err != nil {
		l.WithError(err).Fatal("Unable to initialize tracer.")
	}
	defer func(tc io.Closer) {
		err := tc.Close()
		if err != nil {
			l.WithError(err).Errorf("Unable to close tracer.")
		}
	}(tc)

	config, err := configuration.NewConfigurator(l).GetConfiguration()
	if err != nil {
		l.WithError(err).Errorf("Unable to load service configuration.")
		return
	}
	transport.InitializeRegistry(config)

	rest.CreateService(l, config, ctx, wg, "/ms/cts", transport.InitResource)

	kafka.CreateConsumers(l, ctx, wg, transport.StatusConsumer(consumerGroupId))

	go tasks.Register(transport.NewStateEvaluationTask(l, config, 5000))

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
