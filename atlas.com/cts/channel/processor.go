package channel

import (
	"atlas-cts/rest/requests"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

func GetAll(l logrus.FieldLogger, span opentracing.Span) ([]Model, error) {
	return requests.SliceProvider[attributes, Model](l, span)(requestChannels(), makeChannel)()
}

func makeChannel(data requests.DataBody[attributes]) (Model, error) {
	att := data.Attributes
	return NewChannelBuilder().
		SetWorldId(att.WorldId).
		SetChannelId(att.ChannelId).
		SetCapacity(att.Capacity).
		SetIpAddress(att.IpAddress).
		SetPort(att.Port).
		Build(), nil
}
