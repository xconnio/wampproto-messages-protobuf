package parsers

import (
	"google.golang.org/protobuf/proto"

	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-protobuf/go/gen"
)

type Register struct {
	gen *gen.Register
}

func NewRegisterFields(gen *gen.Register) messages.RegisterFields {
	return &Register{gen: gen}
}

func (r *Register) RequestID() uint64 {
	return r.gen.GetRequestId()
}

func (r *Register) Options() map[string]any {
	options := make(map[string]any)
	if r.gen.Match.String() != "" {
		setDetail(&options, "match", r.gen.Match.String())
	}
	if r.gen.Invoke.String() != "" {
		setDetail(&options, "invoke", r.gen.Invoke.String())
	}
	return options
}

func (r *Register) Procedure() string {
	return r.gen.GetProcedure()
}

func RegisterToProtobuf(register *messages.Register) ([]byte, error) {
	msg := &gen.Register{
		RequestId: register.RequestID(),
		Procedure: register.Procedure(),
	}

	matchString, ok := register.Options()["match"].(string)
	if ok {
		matchValue, ok := gen.Register_Match_value[matchString]
		if ok {
			msg.Match = gen.Register_Match(matchValue)
		}
	}

	invokeString, ok := register.Options()["invoke"].(string)
	if ok {
		invokeValue, ok := gen.Register_Invoke_value[invokeString]
		if ok {
			msg.Invoke = gen.Register_Invoke(invokeValue)
		}
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return PrependHeader(messages.MessageTypeRegister, data, nil), nil
}

func ProtobufToRegister(data []byte) (*messages.Register, error) {
	msg := &gen.Register{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}

	return messages.NewRegisterWithFields(NewRegisterFields(msg)), nil
}
