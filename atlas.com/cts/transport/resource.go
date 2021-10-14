package transport

import (
	"atlas-cts/configuration"
	"atlas-cts/json"
	"atlas-cts/rest"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

const (
	GetTransport  = "get_transport"
	GetTransports = "get_transports"
)

func InitResource(router *mux.Router, l logrus.FieldLogger, c *configuration.Configuration) {
	tr := router.PathPrefix("/transports").Subrouter()
	tr.HandleFunc("/", registerGetTransport(l, c)).Queries("source", "{source}", "destination", "{destination}").Methods(http.MethodGet)
	tr.HandleFunc("/", registerGetTransports(l, c)).Methods(http.MethodGet)
}

func registerGetTransports(l logrus.FieldLogger, c *configuration.Configuration) http.HandlerFunc {
	return rest.RetrieveSpan(GetTransports, func(span opentracing.Span) http.HandlerFunc {
		return handleGetTransports(l, c)(span)
	})
}

func handleGetTransports(l logrus.FieldLogger, c *configuration.Configuration) func(span opentracing.Span) http.HandlerFunc {
	return func(span opentracing.Span) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			transports, err := GetAll(l, c)()
			if err != nil {
				l.WithError(err).Errorf("Unable to get transports configured for service.")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			response := &dataListContainer{}
			for _, t := range transports {
				response.Data = append(response.Data, makeDataBody(t))
			}

			err = json.ToJSON(response, w)
			if err != nil {
				l.WithError(err).Errorf("Writing output")
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	}
}

func registerGetTransport(l logrus.FieldLogger, c *configuration.Configuration) http.HandlerFunc {
	return rest.RetrieveSpan(GetTransport, func(span opentracing.Span) http.HandlerFunc {
		return handleGetTransport(l, c)(span)
	})
}

func handleGetTransport(l logrus.FieldLogger, c *configuration.Configuration) func(span opentracing.Span) http.HandlerFunc {
	return func(span opentracing.Span) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			filters := make([]Filter, 0)
			if val, ok := mux.Vars(r)["source"]; ok {
				source, err := strconv.Atoi(val)
				if err != nil {
					l.WithError(err).Errorf("Unable to properly parse source from path.")
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				filters = append(filters, MatchSource(uint32(source)))
			}
			if val, ok := mux.Vars(r)["destination"]; ok {
				destination, err := strconv.Atoi(val)
				if err != nil {
					l.WithError(err).Errorf("Unable to properly parse destination from path.")
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				filters = append(filters, MatchDestination(uint32(destination)))
			}

			transports, err := GetFiltered(l, c)(filters...)()
			if err != nil {
				l.WithError(err).Errorf("Unable to get transports configured for service.")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			response := &dataListContainer{}
			for _, t := range transports {
				response.Data = append(response.Data, makeDataBody(t))
			}

			err = json.ToJSON(response, w)
			if err != nil {
				l.WithError(err).Errorf("Writing output")
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	}
}

func makeDataBody(t *Model) dataBody {
	return dataBody{
		Id:   "",
		Type: "",
		Attributes: attributes{
			Enabled:     t.Enabled(),
			Source:      t.Source(),
			Departure:   t.Departure(),
			Transport:   t.Transport(),
			Arrival:     t.Arrival(),
			Destination: t.Destination(),
		},
	}
}
