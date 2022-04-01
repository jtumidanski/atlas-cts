package character

import (
	"atlas-cts/kafka"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type changeMapEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	CharacterId uint32 `json:"characterId"`
	MapId       uint32 `json:"mapId"`
	PortalId    uint32 `json:"portalId"`
}

func emitChangeMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, characterId uint32, mapId uint32, portalId uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_CHANGE_MAP_COMMAND")
	return func(worldId byte, channelId byte, characterId uint32, mapId uint32, portalId uint32) {
		event := &changeMapEvent{WorldId: worldId, ChannelId: channelId, CharacterId: characterId, MapId: mapId, PortalId: portalId}
		producer(kafka.CreateKey(int(characterId)), event)
	}
}
