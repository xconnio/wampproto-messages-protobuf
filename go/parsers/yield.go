package parsers

import (
	"google.golang.org/protobuf/proto"

	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-go/serializers"
	"github.com/xconnio/wampproto-protobuf/go/gen"
)

type Yield struct {
	gen *gen.Yield
	ex  *PayloadExpander
}

func NewYieldFields(gen *gen.Yield, payload []byte) messages.YieldFields {
	return &Yield{
		gen: gen,
		ex:  &PayloadExpander{payload: payload, serializer: gen.GetPayloadSerializerId()},
	}
}

func (y *Yield) RequestID() uint64 {
	return y.gen.GetRequestId()
}

func (y *Yield) Options() map[string]any {
	return map[string]any{}
}

func (y *Yield) Args() []any {
	return y.ex.Args()
}

func (y *Yield) KwArgs() map[string]any {
	return y.ex.Kwargs()
}

func (y *Yield) PayloadIsBinary() bool {
	return true
}

func (y *Yield) Payload() []byte {
	return y.ex.Payload()
}

func (y *Yield) PayloadSerializer() uint64 {
	return y.gen.GetPayloadSerializerId()
}

func YieldToProtobuf(yield *messages.Yield) ([]byte, error) {
	msg := &gen.Yield{
		RequestId: yield.RequestID(),
	}

	payloadSerializer := selectPayloadSerializer(yield.Options())
	msg.PayloadSerializerId = payloadSerializer

	var payload []byte
	if yield.PayloadIsBinary() {
		payload = yield.Payload()
	} else {
		var err error
		payload, err = serializers.SerializePayload(payloadSerializer, yield.Args(), yield.KwArgs())
		if err != nil {
			return nil, err
		}
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return PrependHeader(messages.MessageTypeYield, data, payload), nil
}

func ProtobufToYield(data, payload []byte) (*messages.Yield, error) {
	msg := &gen.Yield{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}

	return messages.NewYieldWithFields(NewYieldFields(msg, payload)), nil
}
