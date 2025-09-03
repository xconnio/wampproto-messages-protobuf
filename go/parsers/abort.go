package parsers

import (
	"google.golang.org/protobuf/proto"

	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-go/serializers"
	"github.com/xconnio/wampproto-protobuf/go/gen"
)

type Abort struct {
	gen *gen.Abort
	ex  *PayloadExpander
}

func NewAbortFields(g *gen.Abort, payload []byte) messages.AbortFields {
	return &Abort{
		gen: g,
		ex:  &PayloadExpander{payload: payload, serializer: g.GetPayloadSerializerId()},
	}
}

func (a *Abort) Details() map[string]any {
	return map[string]any{}
}

func (a *Abort) Reason() string {
	return a.gen.Reason
}

func (a *Abort) Args() []any {
	return a.ex.Args()
}

func (a *Abort) KwArgs() map[string]any {
	return a.ex.Kwargs()
}

func AbortToProtobuf(abort *messages.Abort) ([]byte, error) {
	msg := &gen.Abort{
		Reason: abort.Reason(),
	}

	payloadSerializer := selectPayloadSerializer(abort.Details())
	msg.PayloadSerializerId = payloadSerializer
	payload, err := serializers.SerializePayload(payloadSerializer, abort.Args(), abort.KwArgs())
	if err != nil {
		return nil, err
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return PrependHeader(messages.MessageTypeAbort, data, payload), nil
}

func ProtobufToAbort(data, payload []byte) (*messages.Abort, error) {
	msg := &gen.Abort{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}

	return messages.NewAbortWithFields(NewAbortFields(msg, payload)), nil
}
