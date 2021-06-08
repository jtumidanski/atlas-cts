package transport

import (
	"atlas-cts/configuration"
	"github.com/sirupsen/logrus"
)

type ListProvider func() ([]*Model, error)

func GetAll(l logrus.FieldLogger, c *configuration.Configuration) ListProvider {
	return func() ([]*Model, error) {
		results := make([]*Model, 0)
		for _, tc := range c.Transports {
			t := &Model{
				enabled:     tc.Enabled,
				source:      tc.Source,
				departure:   tc.Departure,
				transport:   tc.Transport,
				arrival:     tc.Arrival,
				destination: tc.Destination,
			}
			results = append(results, t)
		}
		return results, nil
	}
}

type Filter func(*Model) bool

func FilterEnabled(value bool) Filter {
	return func(model *Model) bool {
		return model.Enabled() == value
	}
}

func MatchSource(source uint32) Filter {
	return func(model *Model) bool {
		return model.Source() == source
	}
}

func MatchDestination(destination uint32) Filter {
	return func(model *Model) bool {
		return model.Destination() == destination
	}
}

func GetFiltered(l logrus.FieldLogger, c *configuration.Configuration) func(filters ...Filter) ListProvider {
	return func(filters ...Filter) ListProvider {
		return func() ([]*Model, error) {
			transports, err := GetAll(l, c)()
			if err != nil {
				return nil, err
			}

			results := make([]*Model, 0)
			for _, t := range transports {
				ok := true
				for _, filter := range filters {
					if !filter(t) {
						ok = false
					}
				}
				if ok {
					results = append(results, t)
				}
			}
			return results, nil
		}
	}
}
