package parsers

import (
	"google.golang.org/protobuf/proto"

	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-go/serializers"
	"github.com/xconnio/wampproto-protobuf/go/gen"
)

type resultFields struct {
	gen *gen.Result
	ex  *PayloadExpander
}

func NewResultFields(gen *gen.Result, payload []byte) messages.ResultFields {
	return &resultFields{
		gen: gen,
		ex:  &PayloadExpander{payload: payload, serializer: gen.GetPayloadSerializerId()},
	}
}

func (r *resultFields) RequestID() uint64 {
	return r.gen.GetRequestId()
}

func (r *resultFields) Details() map[string]any {
	return map[string]any{}
}

func (r *resultFields) Args() []any {
	return r.ex.Args()
}

func (r *resultFields) KwArgs() map[string]any {
	return r.ex.Kwargs()
}

func (r *resultFields) PayloadIsBinary() bool {
	return true
}

func (r *resultFields) Payload() []byte {
	return r.ex.Payload()
}

func (r *resultFields) PayloadSerializer() uint64 {
	return r.gen.GetPayloadSerializerId()
}

func ResultToProtobuf(result *messages.Result) ([]byte, error) {
	msg := &gen.Yield{
		RequestId: result.RequestID(),
	}

	payloadSerializer := selectPayloadSerializer(result.Details())
	msg.PayloadSerializerId = payloadSerializer

	var payload []byte
	if result.PayloadIsBinary() {
		payload = result.Payload()
	} else {
		var err error
		payload, err = serializers.SerializePayload(payloadSerializer, result.Args(), result.KwArgs())
		if err != nil {
			return nil, err
		}
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return PrependHeader(messages.MessageTypeResult, data, payload), nil
}

func ProtobufToResult(data, payload []byte) (*messages.Result, error) {
	msg := &gen.Result{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}

	return messages.NewResultWithFields(NewResultFields(msg, payload)), nil
}
