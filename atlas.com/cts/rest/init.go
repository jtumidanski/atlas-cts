package rest

import (
	"atlas-cts/configuration"
	"context"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

type RouteInitializer func(*mux.Router, logrus.FieldLogger, *configuration.Configuration)

func CreateService(l *logrus.Logger, c *configuration.Configuration, ctx context.Context, wg *sync.WaitGroup, basePath string, initializers ...RouteInitializer) {
	go NewServer(l, ctx, wg, ProduceRoutes(c, basePath, initializers...))
}

func ProduceRoutes(c *configuration.Configuration, basePath string, initializers ...RouteInitializer) func(l logrus.FieldLogger) http.Handler {
	return func(l logrus.FieldLogger) http.Handler {
		router := mux.NewRouter().PathPrefix(basePath).Subrouter().StrictSlash(true)
		router.Use(CommonHeader)

		for _, initializer := range initializers {
			initializer(router, l, c)
		}

		return router
	}
}