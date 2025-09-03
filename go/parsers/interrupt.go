package parsers

import (
	"google.golang.org/protobuf/proto"

	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-protobuf/go/gen"
)

type Interrupt struct {
	gen *gen.Interrupt
}

func NewInterruptFields(gen *gen.Interrupt) messages.InterruptFields {
	return &Interrupt{gen: gen}
}

func (i *Interrupt) RequestID() uint64 {
	return i.gen.RequestId
}

func (i *Interrupt) Options() map[string]any {
	return map[string]any{}
}

func InterruptToProtobuf(interrupt *messages.Interrupt) ([]byte, error) {
	msg := &gen.Interrupt{
		RequestId: interrupt.RequestID(),
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return PrependHeader(messages.MessageTypeInterrupt, data, nil), nil
}

func ProtobufToInterrupt(data []byte) (*messages.Interrupt, error) {
	msg := &gen.Interrupt{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}

	return messages.NewInterruptWithFields(NewInterruptFields(msg)), nil
}
