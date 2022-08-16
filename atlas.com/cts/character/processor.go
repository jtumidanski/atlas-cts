package character

import (
	"atlas-cts/model"
	"atlas-cts/portal"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

func WarpToPortal(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, characterId uint32, mapId uint32, p model.IdProvider[uint32]) {
	return func(worldId byte, channelId byte, characterId uint32, mapId uint32, p model.IdProvider[uint32]) {
		emitChangeMap(l, span)(worldId, channelId, characterId, mapId, p())
	}
}

func WarpRandom(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, characterId uint32, mapId uint32) {
	return func(worldId byte, channelId byte, characterId uint32, mapId uint32) {
		WarpToPortal(l, span)(worldId, channelId, characterId, mapId, portal.RandomPortalIdProvider(l, span)(mapId))
	}
}
