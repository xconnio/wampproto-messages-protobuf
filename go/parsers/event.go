package parsers

import (
	"google.golang.org/protobuf/proto"

	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-go/serializers"
	"github.com/xconnio/wampproto-go/util"
	"github.com/xconnio/wampproto-protobuf/go/gen"
)

type Event struct {
	gen *gen.Event
	ex  *PayloadExpander
}

func NewEventFields(gen *gen.Event, payload []byte) messages.EventFields {
	return &Event{
		gen: gen,
		ex:  &PayloadExpander{payload: payload, serializer: gen.GetPayloadSerializerId()},
	}
}

func (e *Event) SubscriptionID() uint64 {
	return e.gen.SubscriptionId
}

func (e *Event) PublicationID() uint64 {
	return e.gen.PublicationId
}

func setDetail(details *map[string]any, key string, value any) {
	if *details == nil {
		*details = make(map[string]any)
	}

	(*details)[key] = value
}

func (e *Event) Details() map[string]any {
	var details map[string]any

	if e.gen.Publisher > 0 {
		setDetail(&details, "publisher", e.gen.Publisher)
	}

	if e.gen.PublisherAuthid != "" {
		setDetail(&details, "publisher_authid", e.gen.PublisherAuthid)
	}

	if e.gen.PublisherAuthrole != "" {
		setDetail(&details, "publisher_authrole", e.gen.PublisherAuthrole)
	}

	if e.gen.Topic != "" {
		setDetail(&details, "topic", e.gen.Topic)
	}

	return details
}

func (e *Event) Args() []any {
	return e.ex.Args()
}

func (e *Event) KwArgs() map[string]any {
	return e.ex.Kwargs()
}

func (e *Event) PayloadIsBinary() bool {
	return true
}

func (e *Event) Payload() []byte {
	return e.ex.Payload()
}

func (e *Event) PayloadSerializer() uint64 {
	return e.gen.GetPayloadSerializerId()
}

func EventToProtobuf(event *messages.Event) ([]byte, error) {
	msg := &gen.Event{
		SubscriptionId: event.SubscriptionID(),
		PublicationId:  event.PublicationID(),
	}

	if publisher, ok := util.AsUInt64(event.Details()["publisher"]); ok {
		msg.Publisher = publisher

		authID, ok := util.AsString(event.Details()["publisher_authid"])
		if ok {
			msg.PublisherAuthid = authID
		}

		authRole, ok := util.AsString(event.Details()["publisher_authrole"])
		if ok {
			msg.PublisherAuthrole = authRole
		}

		topic, ok := util.AsString(event.Details()["topic"])
		if ok {
			msg.Topic = topic
		}
	}

	payloadSerializer := selectPayloadSerializer(event.Details())
	msg.PayloadSerializerId = payloadSerializer
	payload, err := serializers.SerializePayload(payloadSerializer, event.Args(), event.KwArgs())
	if err != nil {
		return nil, err
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return PrependHeader(messages.MessageTypeEvent, data, payload), nil
}

func ProtobufToEvent(data []byte, payload []byte) (*messages.Event, error) {
	msg := &gen.Event{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}

	return messages.NewEventWithFields(NewEventFields(msg, payload)), nil
}
