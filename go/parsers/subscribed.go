package parsers

import (
	"google.golang.org/protobuf/proto"

	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-protobuf/go/gen"
)

type Subscribed struct {
	gen *gen.Subscribed
}

func NewSubscribedFields(gen *gen.Subscribed) messages.SubscribedFields {
	return &Subscribed{gen: gen}
}

func (s *Subscribed) RequestID() uint64 {
	return s.gen.GetRequestId()
}

func (s *Subscribed) SubscriptionID() uint64 {
	return s.gen.SubscriptionId
}

func SubscribedToProtobuf(subscribed *messages.Subscribed) ([]byte, error) {
	msg := &gen.Subscribed{
		RequestId:      subscribed.RequestID(),
		SubscriptionId: subscribed.SubscriptionID(),
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return PrependHeader(messages.MessageTypeSubscribed, data, nil), nil
}

func ProtobufToSubscribed(data []byte) (*messages.Subscribed, error) {
	msg := &gen.Subscribed{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}

	return messages.NewSubscribedWithFields(NewSubscribedFields(msg)), nil
}
