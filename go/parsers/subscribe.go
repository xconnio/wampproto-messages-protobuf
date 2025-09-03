package parsers

import (
	"google.golang.org/protobuf/proto"

	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-protobuf/go/gen"
)

type Subscribe struct {
	gen *gen.Subscribe
}

func NewSubscribeFields(gen *gen.Subscribe) messages.SubscribeFields {
	return &Subscribe{gen: gen}
}

func (s *Subscribe) RequestID() uint64 {
	return s.gen.GetRequestId()
}

func (s *Subscribe) Options() map[string]any {
	options := make(map[string]any)
	if s.gen.Match.String() != "" {
		setDetail(&options, "match", s.gen.Match.String())
	}
	return options
}

func (s *Subscribe) Topic() string {
	return s.gen.GetTopic()
}

func SubscribeToProtobuf(subscribe *messages.Subscribe) ([]byte, error) {
	msg := &gen.Subscribe{
		RequestId: subscribe.RequestID(),
		Topic:     subscribe.Topic(),
	}

	matchString, ok := subscribe.Options()["match"].(string)
	if ok {
		matchValue, ok := gen.Subscribe_Match_value[matchString]
		if ok {
			msg.Match = gen.Subscribe_Match(matchValue)
		}
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return PrependHeader(messages.MessageTypeSubscribe, data, nil), nil
}

func ProtobufToSubscribe(data []byte) (*messages.Subscribe, error) {
	msg := &gen.Subscribe{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}

	return messages.NewSubscribeWithFields(NewSubscribeFields(msg)), nil
}
