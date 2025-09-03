package parsers

import (
	"google.golang.org/protobuf/proto"

	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-protobuf/go/gen"
)

type Cancel struct {
	gen *gen.Cancel
}

func NewCancelFields(gen *gen.Cancel) messages.CancelFields {
	return &Cancel{gen: gen}
}

func (c *Cancel) RequestID() uint64 {
	return c.gen.RequestId
}

func (c *Cancel) Options() map[string]any {
	return map[string]any{}
}

func CancelToProtobuf(cancel *messages.Cancel) ([]byte, error) {
	msg := &gen.Cancel{
		RequestId: cancel.RequestID(),
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return PrependHeader(messages.MessageTypeCancel, data, nil), nil
}

func ProtobufToCancel(data []byte) (*messages.Cancel, error) {
	msg := &gen.Cancel{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}

	return messages.NewCancelWithFields(NewCancelFields(msg)), nil
}
