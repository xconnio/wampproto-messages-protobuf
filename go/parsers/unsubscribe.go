package parsers

import (
	"google.golang.org/protobuf/proto"

	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-protobuf/go/gen"
)

type Unsubscribe struct {
	gen *gen.UnSubscribe
}

func NewUnsubscribeFields(gen *gen.UnSubscribe) messages.UnsubscribeFields {
	return &Unsubscribe{gen: gen}
}

func (u *Unsubscribe) RequestID() uint64 {
	return u.gen.GetRequestId()
}

func (u *Unsubscribe) SubscriptionID() uint64 {
	return u.gen.GetSubscriptionId()
}

func UnsubscribeToProtobuf(unsubscribe *messages.Unsubscribe) ([]byte, error) {
	msg := &gen.UnSubscribe{
		RequestId:      unsubscribe.RequestID(),
		SubscriptionId: unsubscribe.SubscriptionID(),
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	byteValue := byte(messages.MessageTypeUnsubscribe & 0xFF)
	return append([]byte{byteValue}, data...), nil
}

func ProtobufToUnsubscribe(data []byte) (*messages.Unsubscribe, error) {
	msg := &gen.UnSubscribe{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}

	return messages.NewUnsubscribeWithFields(NewUnsubscribeFields(msg)), nil
}
