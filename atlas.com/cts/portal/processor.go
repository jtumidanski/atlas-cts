package portal

import (
	"atlas-cts/model"
	"atlas-cts/rest/requests"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

func RandomPortalProvider(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32) model.Provider[Model] {
	return func(mapId uint32) model.Provider[Model] {
		return model.SliceProviderToProviderAdapter[Model](InMapProvider(l, span)(mapId), model.RandomPreciselyOneFilter[Model])
	}
}

func getId(m Model) (uint32, error) {
	return m.Id(), nil
}

func RandomPortalIdProvider(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32) model.IdProvider[uint32] {
	return func(mapId uint32) model.IdProvider[uint32] {
		return model.ProviderToIdProviderAdapter[Model, uint32](RandomPortalProvider(l, span)(mapId), getId)
	}
}

func InMapProvider(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32) model.SliceProvider[Model] {
	return func(mapId uint32) model.SliceProvider[Model] {
		return requests.SliceProvider[attributes, Model](l, span)(requestAll(mapId), makePortal)
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
