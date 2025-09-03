package parsers

import (
	"google.golang.org/protobuf/proto"

	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-protobuf/go/gen"
)

type Unregistered struct {
	gen *gen.Unregistered
}

func NewUnregisteredFields(gen *gen.Unregistered) messages.UnregisteredFields {
	return &Unregistered{gen: gen}
}

func (u *Unregistered) RequestID() uint64 {
	return u.gen.GetRequestId()
}

func UnregisteredToProtobuf(unregistered *messages.Unregistered) ([]byte, error) {
	msg := &gen.Unregistered{
		RequestId: unregistered.RequestID(),
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return PrependHeader(messages.MessageTypeUnregistered, data, nil), nil
}

func ProtobufToUnregistered(data []byte) (*messages.Unregistered, error) {
	msg := &gen.Unregistered{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}

	return messages.NewUnregisteredWithFields(NewUnregisteredFields(msg)), nil
}
