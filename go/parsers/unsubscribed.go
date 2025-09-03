package parsers

import (
	"google.golang.org/protobuf/proto"

	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-protobuf/go/gen"
)

type Unsubscribed struct {
	gen *gen.Unsubscribed
}

func NewUnsubscribedFields(gen *gen.Unsubscribed) messages.UnsubscribedFields {
	return &Unsubscribed{gen: gen}
}

func (u *Unsubscribed) RequestID() uint64 {
	return u.gen.GetRequestId()
}

func UnsubscribedToProtobuf(unsubscribed *messages.Unsubscribed) ([]byte, error) {
	msg := &gen.Unsubscribed{
		RequestId: unsubscribed.RequestID(),
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return PrependHeader(messages.MessageTypeUnsubscribed, data, nil), nil
}

func ProtobufToUnsubscribed(data []byte) (*messages.Unsubscribed, error) {
	msg := &gen.Unsubscribed{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}

	return messages.NewUnsubscribedWithFields(NewUnsubscribedFields(msg)), nil
}
