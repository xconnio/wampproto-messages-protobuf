package parsers

import (
	"google.golang.org/protobuf/proto"

	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-protobuf/go/gen"
)

type Unsubscribe struct {
	gen *gen.Unsubscribe
}

func NewUnsubscribeFields(gen *gen.Unsubscribe) messages.UnsubscribeFields {
	return &Unsubscribe{gen: gen}
}

func (u *Unsubscribe) RequestID() uint64 {
	return u.gen.GetRequestId()
}

func (u *Unsubscribe) SubscriptionID() uint64 {
	return u.gen.GetSubscriptionId()
}

func UnsubscribeToProtobuf(unsubscribe *messages.Unsubscribe) ([]byte, error) {
	msg := &gen.Unsubscribe{
		RequestId:      unsubscribe.RequestID(),
		SubscriptionId: unsubscribe.SubscriptionID(),
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return PrependHeader(messages.MessageTypeUnsubscribe, data, nil), nil
}

func ProtobufToUnsubscribe(data []byte) (*messages.Unsubscribe, error) {
	msg := &gen.Unsubscribe{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}

	return messages.NewUnsubscribeWithFields(NewUnsubscribeFields(msg)), nil
}
