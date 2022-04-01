package transport

import (
	"atlas-cts/channel"
	"atlas-cts/character"
	"atlas-cts/configuration"
	_map "atlas-cts/map"
	"atlas-cts/model"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"math"
	"time"
)

func InitializeRegistry(c *configuration.Configuration) {
	for _, tc := range c.Transports {
		t := Model{
			enabled:     tc.Enabled,
			source:      tc.Source,
			departure:   tc.Departure,
			transport:   tc.Transport,
			arrival:     tc.Arrival,
			destination: tc.Destination,
			state:       "IN_PROGRESS",
		}
		GetRegistry().Add(t)
	}
}

func AllModelProvider(_ logrus.FieldLogger) model.SliceProvider[Model] {
	return func() ([]Model, error) {
		return GetRegistry().GetAll(), nil
	}
}

func FilterEnabled(value bool) model.Filter[Model] {
	return func(model Model) bool {
		return model.Enabled() == value
	}
}

func MatchSource(source uint32) model.Filter[Model] {
	return func(model Model) bool {
		return model.Source() == source
	}
}

func MatchDestination(destination uint32) model.Filter[Model] {
	return func(model Model) bool {
		return model.Destination() == destination
	}
}

const (
	StateBoarding   = "BOARDING"
	StatePreparing  = "PREPARING"
	StateInProgress = "IN_PROGRESS"
)

func getState(t time.Time, c configuration.TransportConfiguration) string {
	ctm := c.RideDuration + c.OpenGateDuration + c.ClosedGateDuration
	tm := uint32(((((t.Hour() * 60) + t.Minute()) * 60) + t.Second()) * 1000)
	cm := uint32(math.Mod(float64(tm), float64(ctm)))
	if cm < c.OpenGateDuration {
		return StateBoarding
	} else if cm < c.OpenGateDuration+c.ClosedGateDuration {
		return StatePreparing
	} else {
		return StateInProgress
	}
}

func UpdateState(l logrus.FieldLogger, span opentracing.Span) func(sourceId uint32, destinationId uint32, newState string) error {
	return func(sourceId uint32, destinationId uint32, newState string) error {
		t, err := GetRegistry().Get(sourceId, destinationId)
		if err != nil {
			l.WithError(err).Errorf("Unable to lookup transport from %d to %d for state change.", sourceId, destinationId)
			return err
		}
		nt := t.updateState(newState)
		err = GetRegistry().Update(nt)
		if err != nil {
			l.WithError(err).Errorf("Unable to update state of transport from %d to %d.", sourceId, destinationId)
			return err
		}
		emitStatusEvent(l, span)(sourceId, destinationId, newState)
		return nil
	}
}

func WarpAllToArrival(l logrus.FieldLogger, span opentracing.Span) func(sourceId uint32, destinationId uint32) {
	return func(sourceId uint32, destinationId uint32) {
		channels, err := channel.GetAll(l, span)
		if err != nil {
			l.WithError(err).Errorf("Unable to locate channels to serve.")
			return
		}
		t, err := GetRegistry().Get(sourceId, destinationId)
		if err != nil {
			l.WithError(err).Errorf("Unable to locate transport from %d to %d.", sourceId, destinationId)
			return
		}

		for _, c := range channels {
			for _, tm := range t.transport {
				go warpAllInMap(l, span)(c.WorldId(), c.ChannelId(), tm, t.Arrival())
			}
		}
	}
}

func WarpAllToTransport(l logrus.FieldLogger, span opentracing.Span) func(sourceId uint32, destinationId uint32) {
	return func(sourceId uint32, destinationId uint32) {
		channels, err := channel.GetAll(l, span)
		if err != nil {
			l.WithError(err).Errorf("Unable to locate channels to serve.")
			return
		}
		t, err := GetRegistry().Get(sourceId, destinationId)
		if err != nil {
			l.WithError(err).Errorf("Unable to locate transport from %d to %d.", sourceId, destinationId)
			return
		}

		for _, c := range channels {
			go warpAllInMap(l, span)(c.WorldId(), c.ChannelId(), t.Departure(), t.Transport()[0])
		}
	}
}

func warpAllInMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, toId uint32) {
	return func(worldId byte, channelId byte, mapId uint32, toId uint32) {
		ids, err := _map.GetCharacterIdsInMap(l, span)(worldId, channelId, mapId)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve characters in world %d, channel %d, map %d for transport.", worldId, channelId, mapId)
			return
		}
		for _, id := range ids {
			character.WarpRandom(l, span)(worldId, channelId, id, toId)
		}
	}
}
