package parsers

import (
	"google.golang.org/protobuf/proto"

	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-go/serializers"
	"github.com/xconnio/wampproto-protobuf/go/gen"
)

type Invocation struct {
	gen *gen.Invocation
	ex  *PayloadExpander
}

func NewInvocationFields(gen *gen.Invocation, payload []byte) messages.InvocationFields {
	return &Invocation{
		gen: gen,
		ex:  &PayloadExpander{payload: payload, serializer: gen.GetPayloadSerializerId()},
	}
}

func (i *Invocation) RequestID() uint64 {
	return i.gen.GetRequestId()
}

func (i *Invocation) RegistrationID() uint64 {
	return i.gen.GetRegistrationId()
}

func (i *Invocation) Details() map[string]any {
	var details map[string]any

	if i.gen.Caller > 0 {
		setDetail(&details, "caller", i.gen.Caller)
	}

	if i.gen.CallerAuthid != "" {
		setDetail(&details, "caller_authid", i.gen.CallerAuthid)
	}

	if i.gen.CallerAuthrole != "" {
		setDetail(&details, "caller_authrole", i.gen.CallerAuthrole)
	}

	if i.gen.Procedure != "" {
		setDetail(&details, "procedure", i.gen.Procedure)
	}

	return details
}

func (i *Invocation) Args() []any {
	return i.ex.Args()
}

func (i *Invocation) KwArgs() map[string]any {
	return i.ex.Kwargs()
}

func (i *Invocation) PayloadIsBinary() bool {
	return true
}

func (i *Invocation) Payload() []byte {
	return i.ex.Payload()
}

func (i *Invocation) PayloadSerializer() uint64 {
	return i.gen.GetPayloadSerializerId()
}

func InvocationToProtobuf(invocation *messages.Invocation) ([]byte, error) {
	msg := &gen.Invocation{
		RequestId:      invocation.RequestID(),
		RegistrationId: invocation.RegistrationID(),
	}

	payloadSerializer := selectPayloadSerializer(invocation.Details())
	msg.PayloadSerializerId = payloadSerializer

	var payload []byte
	if invocation.PayloadIsBinary() {
		payload = invocation.Payload()
	} else {
		var err error
		payload, err = serializers.SerializePayload(payloadSerializer, invocation.Args(), invocation.KwArgs())
		if err != nil {
			return nil, err
		}
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return PrependHeader(messages.MessageTypeInvocation, data, payload), nil
}

func ProtobufToInvocation(data, payload []byte) (*messages.Invocation, error) {
	msg := &gen.Invocation{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}

	return messages.NewInvocationWithFields(NewInvocationFields(msg, payload)), nil
}
