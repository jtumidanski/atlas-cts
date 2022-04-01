package transport

import (
	"atlas-cts/kafka"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type statusEvent struct {
	SourceId      uint32
	DestinationId uint32
	State         string
}

func emitStatusEvent(l logrus.FieldLogger, span opentracing.Span) func(sourceId uint32, destinationId uint32, state string) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_TRANSPORT_STATUS_EVENT")
	return func(sourceId uint32, destinationId uint32, state string) {
		e := &statusEvent{SourceId: sourceId, DestinationId: destinationId, State: state}
		producer(kafka.CreateKey(int(sourceId)), e)
	}
}
