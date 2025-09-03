package parsers

import (
	"google.golang.org/protobuf/proto"

	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-protobuf/go/gen"
)

type GoodBye struct {
	gen *gen.Goodbye
}

func NewGoodByeFields(gen *gen.Goodbye) messages.GoodByeFields {
	return &GoodBye{gen: gen}
}

func (g *GoodBye) Details() map[string]any {
	return map[string]any{}
}

func (g *GoodBye) Reason() string {
	return g.gen.Reason
}

func GoodbyeToProtobuf(goodbye *messages.GoodBye) ([]byte, error) {
	msg := &gen.Goodbye{
		Reason: goodbye.Reason(),
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return PrependHeader(messages.MessageTypeGoodbye, data, nil), nil
}

func ProtobufToGoodbye(data []byte) (*messages.GoodBye, error) {
	msg := &gen.Goodbye{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}

	return messages.NewGoodByeWithFields(NewGoodByeFields(msg)), nil
}
