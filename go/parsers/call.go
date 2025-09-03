package parsers

import (
	"google.golang.org/protobuf/proto"

	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-go/serializers"
	"github.com/xconnio/wampproto-protobuf/go/gen"
)

type Call struct {
	gen *gen.Call
	ex  *PayloadExpander
}

func NewCallFields(callGen *gen.Call, payload []byte) messages.CallFields {
	return &Call{
		gen: callGen,
		ex:  &PayloadExpander{payload: payload, serializer: callGen.GetPayloadSerializerId()},
	}
}

func (c *Call) RequestID() uint64 {
	return c.gen.GetRequestId()
}

func (c *Call) Options() map[string]any {
	return map[string]any{}
}

func (c *Call) Procedure() string {
	return c.gen.GetProcedure()
}

func (c *Call) Args() []any {
	return c.ex.Args()
}

func (c *Call) KwArgs() map[string]any {
	return c.ex.Kwargs()
}

func (c *Call) PayloadIsBinary() bool {
	return true
}

func (c *Call) Payload() []byte {
	return c.ex.Payload()
}

func (c *Call) PayloadSerializer() uint64 {
	return c.gen.GetPayloadSerializerId()
}

func CallToProtobuf(call *messages.Call) ([]byte, error) {
	msg := &gen.Call{
		RequestId: call.RequestID(),
		Procedure: call.Procedure(),
	}

	payloadSerializer := selectPayloadSerializer(call.Options())
	msg.PayloadSerializerId = payloadSerializer

	var payload []byte
	if call.PayloadIsBinary() {
		payload = call.Payload()
	} else {
		var err error
		payload, err = serializers.SerializePayload(payloadSerializer, call.Args(), call.KwArgs())
		if err != nil {
			return nil, err
		}
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return PrependHeader(messages.MessageTypeCall, data, payload), nil
}

func ProtobufToCall(data, payload []byte) (*messages.Call, error) {
	msg := &gen.Call{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}

	return messages.NewCallWithFields(NewCallFields(msg, payload)), nil
}
