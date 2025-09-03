package parsers

import (
	"google.golang.org/protobuf/proto"

	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-go/serializers"
	"github.com/xconnio/wampproto-protobuf/go/gen"
)

type Error struct {
	gen *gen.Error
	ex  *PayloadExpander
}

func NewErrorFields(gen *gen.Error, payload []byte) messages.ErrorFields {
	return &Error{
		gen: gen,
		ex:  &PayloadExpander{payload: payload, serializer: gen.GetPayloadSerializerId()},
	}
}

func (e *Error) MessageType() uint64 {
	return e.gen.MessageType
}

func (e *Error) RequestID() uint64 {
	return e.gen.RequestId
}

func (e *Error) Details() map[string]any {
	return map[string]any{}
}

func (e *Error) URI() string {
	return e.gen.Uri
}

func (e *Error) Args() []any {
	return e.ex.Args()
}

func (e *Error) KwArgs() map[string]any {
	return e.ex.Kwargs()
}

func (e *Error) PayloadIsBinary() bool {
	return true
}

func (e *Error) Payload() []byte {
	return e.ex.Payload()
}

func (e *Error) PayloadSerializer() uint64 {
	return e.gen.GetPayloadSerializerId()
}

func ErrorToProtobuf(errMsg *messages.Error) ([]byte, error) {
	msg := &gen.Error{
		MessageType: errMsg.MessageType(),
		RequestId:   errMsg.RequestID(),
		Uri:         errMsg.URI(),
	}

	payloadSerializer := selectPayloadSerializer(errMsg.Details())
	msg.PayloadSerializerId = payloadSerializer

	var payload []byte
	if errMsg.PayloadIsBinary() {
		payload = errMsg.Payload()
	} else {
		var err error
		payload, err = serializers.SerializePayload(payloadSerializer, errMsg.Args(), errMsg.KwArgs())
		if err != nil {
			return nil, err
		}
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return PrependHeader(messages.MessageTypeError, data, payload), nil
}

func ProtobufToError(data, payload []byte) (*messages.Error, error) {
	msg := &gen.Error{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}

	return messages.NewErrorWithFields(NewErrorFields(msg, payload)), nil
}
