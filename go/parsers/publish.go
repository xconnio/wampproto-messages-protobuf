package parsers

import (
	"google.golang.org/protobuf/proto"

	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-go/serializers"
	"github.com/xconnio/wampproto-protobuf/go/gen"
)

type Publish struct {
	gen *gen.Publish
	ex  *PayloadExpander
}

func NewPublishFields(gen *gen.Publish, payload []byte) messages.PublishFields {
	return &Publish{
		gen: gen,
		ex:  &PayloadExpander{payload: payload, serializer: gen.GetPayloadSerializerId()},
	}
}

func (p *Publish) RequestID() uint64 {
	return p.gen.GetRequestId()
}

func (p *Publish) Options() map[string]any {
	var details map[string]any

	if !p.gen.ExcludeMe {
		setDetail(&details, "exclude_me", false)
	}

	if p.gen.Acknowledge {
		setDetail(&details, "acknowledge", true)
	}

	return details
}

func (p *Publish) Topic() string {
	return p.gen.GetTopic()
}

func (p *Publish) Args() []any {
	return p.ex.Args()
}

func (p *Publish) KwArgs() map[string]any {
	return p.ex.Kwargs()
}

func (p *Publish) PayloadIsBinary() bool {
	return true
}

func (p *Publish) Payload() []byte {
	return p.ex.Payload()
}

func (p *Publish) PayloadSerializer() uint64 {
	return p.gen.GetPayloadSerializerId()
}

func PublishToProtobuf(publish *messages.Publish) ([]byte, error) {
	msg := &gen.Publish{
		RequestId: publish.RequestID(),
		Topic:     publish.Topic(),
	}
	acknowledge, ok := publish.Options()["acknowledge"].(bool)
	if ok {
		msg.Acknowledge = acknowledge
	}

	excludeMe, ok := publish.Options()["exclude_me"].(bool)
	if ok {
		msg.ExcludeMe = excludeMe
	}

	payloadSerializer := selectPayloadSerializer(publish.Options())
	msg.PayloadSerializerId = payloadSerializer

	var payload []byte
	if publish.PayloadIsBinary() {
		payload = publish.Payload()
	} else {
		var err error
		payload, err = serializers.SerializePayload(payloadSerializer, publish.Args(), publish.KwArgs())
		if err != nil {
			return nil, err
		}
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return PrependHeader(messages.MessageTypePublish, data, payload), nil
}

func ProtobufToPublish(data, payload []byte) (*messages.Publish, error) {
	msg := &gen.Publish{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}

	return messages.NewPublishWithFields(NewPublishFields(msg, payload)), nil
}
