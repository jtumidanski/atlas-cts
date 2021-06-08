package transport

import (
	"atlas-cts/configuration"
	"atlas-cts/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func HandleGetTransports(l logrus.FieldLogger, c *configuration.Configuration) http.HandlerFunc {
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

func HandleGetTransport(l logrus.FieldLogger, c *configuration.Configuration) http.HandlerFunc {
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

func makeDataBody(t *Model) dataBody {
	return dataBody{
		Id:         "",
		Type:       "",
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
