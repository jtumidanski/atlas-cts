package transport

import (
	"atlas-cts/kafka"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	consumerNameStatus = "transport_status_event"
	topicNameStatus    = "TOPIC_TRANSPORT_STATUS_EVENT"
)

func StatusConsumer(groupId string) kafka.ConsumerConfig {
	return kafka.NewConsumerConfig[statusEvent](consumerNameStatus, topicNameStatus, groupId, handleStatus())
}

func handleStatus() kafka.HandlerFunc[statusEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event statusEvent) {
		if event.State == StateBoarding {
			WarpAllToArrival(l, span)(event.SourceId, event.DestinationId)
		} else if event.State == StateInProgress {
			WarpAllToTransport(l, span)(event.SourceId, event.DestinationId)
		}
	}
}
