package parsers

import (
	"google.golang.org/protobuf/proto"

	"github.com/xconnio/wampproto-go/auth"
	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-protobuf/go/gen"
)

type Challenge struct {
	gen *gen.Challenge
}

func NewChallengeFields(gen *gen.Challenge) messages.ChallengeFields {
	return &Challenge{gen: gen}
}

func (c *Challenge) AuthMethod() string {
	return c.gen.Authmethod
}

func (c *Challenge) Extra() map[string]any {
	extra := make(map[string]any)
	extra["challenge"] = c.gen.Challenge

	if c.gen.Salt != "" {
		extra["salt"] = c.gen.Salt
		extra["iterations"] = c.gen.Iterations
		extra["keylen"] = c.gen.Keylen
	}
	return extra
}

func ChallengeToProtobuf(challenge *messages.Challenge) ([]byte, error) {
	challengeValue, _ := challenge.Extra()["challenge"].(string)

	msg := &gen.Challenge{
		Authmethod: challenge.AuthMethod(),
		Challenge:  challengeValue,
	}

	if challenge.AuthMethod() == auth.MethodCRA {
		saltString, ok := challenge.Extra()["salt"].(string)
		if ok {
			msg.Salt = saltString
		}

		iterations, ok := challenge.Extra()["iterations"].(uint32)
		if ok {
			msg.Iterations = iterations
		}

		keylen, ok := challenge.Extra()["keylen"].(uint32)
		if ok {
			msg.Keylen = keylen
		}
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return PrependHeader(messages.MessageTypeChallenge, data, nil), nil
}

func ProtobufToChallenge(data []byte) (*messages.Challenge, error) {
	msg := &gen.Challenge{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}

	return messages.NewChallengeWithFields(NewChallengeFields(msg)), nil
}
