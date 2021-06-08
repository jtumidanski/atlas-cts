package rest

import (
	"atlas-cts/configuration"
	"atlas-cts/transport"
	"context"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

func CreateRestService(l *logrus.Logger, c *configuration.Configuration, ctx context.Context, wg *sync.WaitGroup) {
	go NewServer(l, ctx, wg, ProduceRoutes(c))
}

func ProduceRoutes(c *configuration.Configuration) func(logrus.FieldLogger) http.Handler {
	return func(l logrus.FieldLogger) http.Handler {
		router := mux.NewRouter().PathPrefix("/ms/cts").Subrouter().StrictSlash(true)
		router.Use(CommonHeader)

		tr := router.PathPrefix("/transports").Subrouter()
		tr.HandleFunc("/", transport.HandleGetTransport(l, c)).Queries("source", "{source}", "destination", "{destination}").Methods(http.MethodGet)
		tr.HandleFunc("/", transport.HandleGetTransports(l, c)).Methods(http.MethodGet)
		return router
	}
}
