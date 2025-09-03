package parsers

import (
	"google.golang.org/protobuf/proto"

	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-protobuf/go/gen"
)

type Authenticate struct {
	gen *gen.Authenticate
}

func NewAuthenticateFields(gen *gen.Authenticate) messages.AuthenticateFields {
	return &Authenticate{gen: gen}
}

func (a *Authenticate) Signature() string {
	return a.gen.Signature
}

func (a *Authenticate) Extra() map[string]any {
	return map[string]any{}
}

func AuthenticateToProtobuf(authenticate *messages.Authenticate) ([]byte, error) {
	msg := &gen.Authenticate{
		Signature: authenticate.Signature(),
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return PrependHeader(messages.MessageTypeAuthenticate, data, nil), nil
}

func ProtobufToAuthenticate(data []byte) (*messages.Authenticate, error) {
	msg := &gen.Authenticate{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}

	return messages.NewAuthenticateWithFields(NewAuthenticateFields(msg)), nil
}
