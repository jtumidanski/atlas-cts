package _map

import (
	"atlas-cts/model"
	"atlas-cts/rest/requests"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

func CharacterIdsInMapModelProvider(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32) model.SliceProvider[uint32] {
	return func(worldId byte, channelId byte, mapId uint32) model.SliceProvider[uint32] {
		return requests.SliceProvider[characterAttributes, uint32](l, span)(requestCharactersInMap(worldId, channelId, mapId), getCharacterId)
	}
}

func GetCharacterIdsInMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32) ([]uint32, error) {
	return func(worldId byte, channelId byte, mapId uint32) ([]uint32, error) {
		return CharacterIdsInMapModelProvider(l, span)(worldId, channelId, mapId)()
	}
}

func getCharacterId(body requests.DataBody[characterAttributes]) (uint32, error) {
	id, err := strconv.ParseUint(body.Id, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(id), nil
}
