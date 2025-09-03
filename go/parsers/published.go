package parsers

import (
	"google.golang.org/protobuf/proto"

	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-protobuf/go/gen"
)

type Published struct {
	gen *gen.Published
}

func NewPublishedFields(gen *gen.Published) messages.PublishedFields {
	return &Published{gen: gen}
}

func (p *Published) RequestID() uint64 {
	return p.gen.GetRequestId()
}

func (p *Published) PublicationID() uint64 {
	return p.gen.GetPublicationId()
}

func PublishedToProtobuf(published *messages.Published) ([]byte, error) {
	msg := &gen.Published{
		RequestId:     published.RequestID(),
		PublicationId: published.PublicationID(),
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return PrependHeader(messages.MessageTypePublished, data, nil), nil
}

func ProtobufToPublished(data []byte) (*messages.Published, error) {
	msg := &gen.Published{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}

	return messages.NewPublishedWithFields(NewPublishedFields(msg)), nil
}
