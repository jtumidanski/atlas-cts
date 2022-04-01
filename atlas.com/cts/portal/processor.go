package portal

import (
	"atlas-cts/rest/requests"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
)

type IdProvider func() uint32

func RandomPortalIdProvider(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32) IdProvider {
	return func(mapId uint32) IdProvider {
		return func() uint32 {
			ps, err := ForMap(l, span)(mapId)
			if err != nil {
				l.WithError(err).Errorf("Unable to retrieve portals for map %d. Defaulting to 0.", mapId)
				return 0
			}
			if len(ps) == 0 {
				l.Warnf("No portals in map %d. Defaulting to zero.", mapId)
				return 0
			}
			return ps[rand.Intn(len(ps))].Id()
		}
	}
}

func ForMap(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32) ([]Model, error) {
	return func(mapId uint32) ([]Model, error) {
		return requests.SliceProvider[attributes, Model](l, span)(requestAll(mapId), makePortal)()
	}
}

func makePortal(body requests.DataBody[attributes]) (Model, error) {
	id, err := strconv.ParseUint(body.Id, 10, 32)
	if err != nil {
		return Model{}, err
	}
	attr := body.Attributes
	return NewPortalModel(uint32(id), attr.Name, attr.Target, attr.TargetMapId, attr.Type, attr.X, attr.Y, attr.ScriptName), nil
}
