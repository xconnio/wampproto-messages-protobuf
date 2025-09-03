package wampprotobuf

import (
	"fmt"

	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-go/serializers"
	"github.com/xconnio/wampproto-protobuf/go/parsers"
)

type ProtobufSerializer struct{}

var _ serializers.Serializer = &ProtobufSerializer{}

func (p *ProtobufSerializer) Serialize(message messages.Message) ([]byte, error) {
	switch message.Type() {
	case messages.MessageTypeHello:
		msg := message.(*messages.Hello)
		return parsers.HelloToProtobuf(msg)
	case messages.MessageTypeWelcome:
		msg := message.(*messages.Welcome)
		return parsers.WelcomeToProtobuf(msg)
	case messages.MessageTypeChallenge:
		msg := message.(*messages.Challenge)
		return parsers.ChallengeToProtobuf(msg)
	case messages.MessageTypeAuthenticate:
		msg := message.(*messages.Authenticate)
		return parsers.AuthenticateToProtobuf(msg)
	case messages.MessageTypeAbort:
		msg := message.(*messages.Abort)
		return parsers.AbortToProtobuf(msg)
	case messages.MessageTypeGoodbye:
		msg := message.(*messages.GoodBye)
		return parsers.GoodbyeToProtobuf(msg)
	case messages.MessageTypeRegister:
		msg := message.(*messages.Register)
		return parsers.RegisterToProtobuf(msg)
	case messages.MessageTypeRegistered:
		msg := message.(*messages.Registered)
		return parsers.RegisteredToProtobuf(msg)
	case messages.MessageTypeUnregister:
		msg := message.(*messages.Unregister)
		return parsers.UnregisterToProtobuf(msg)
	case messages.MessageTypeUnregistered:
		msg := message.(*messages.Unregistered)
		return parsers.UnregisteredToProtobuf(msg)
	case messages.MessageTypePublish:
		msg := message.(*messages.Publish)
		return parsers.PublishToProtobuf(msg)
	case messages.MessageTypePublished:
		msg := message.(*messages.Published)
		return parsers.PublishedToProtobuf(msg)
	case messages.MessageTypeEvent:
		msg := message.(*messages.Event)
		return parsers.EventToProtobuf(msg)
	case messages.MessageTypeCall:
		msg := message.(*messages.Call)
		return parsers.CallToProtobuf(msg)
	case messages.MessageTypeInvocation:
		msg := message.(*messages.Invocation)
		return parsers.InvocationToProtobuf(msg)
	case messages.MessageTypeYield:
		msg := message.(*messages.Yield)
		return parsers.YieldToProtobuf(msg)
	case messages.MessageTypeResult:
		msg := message.(*messages.Result)
		return parsers.ResultToProtobuf(msg)
	case messages.MessageTypeSubscribe:
		msg := message.(*messages.Subscribe)
		return parsers.SubscribeToProtobuf(msg)
	case messages.MessageTypeSubscribed:
		msg := message.(*messages.Subscribed)
		return parsers.SubscribedToProtobuf(msg)
	case messages.MessageTypeUnsubscribe:
		msg := message.(*messages.Unsubscribe)
		return parsers.UnsubscribeToProtobuf(msg)
	case messages.MessageTypeUnsubscribed:
		msg := message.(*messages.Unsubscribed)
		return parsers.UnsubscribedToProtobuf(msg)
	case messages.MessageTypeCancel:
		msg := message.(*messages.Cancel)
		return parsers.CancelToProtobuf(msg)
	case messages.MessageTypeInterrupt:
		msg := message.(*messages.Interrupt)
		return parsers.InterruptToProtobuf(msg)
	case messages.MessageTypeError:
		msg := message.(*messages.Error)
		return parsers.ErrorToProtobuf(msg)
	default:
		return nil, fmt.Errorf("unknown message type: %v", message.Type())
	}
}

func (p *ProtobufSerializer) Deserialize(bytes []byte) (messages.Message, error) {
	messageData, payloadData, err := parsers.ExtractMessage(bytes)
	if err != nil {
		return nil, err
	}

	switch uint64(bytes[0]) {
	case messages.MessageTypeHello:
		return parsers.ProtobufToHello(messageData)
	case messages.MessageTypeWelcome:
		return parsers.ProtobufToWelcome(messageData)
	case messages.MessageTypeChallenge:
		return parsers.ProtobufToChallenge(messageData)
	case messages.MessageTypeAuthenticate:
		return parsers.ProtobufToAuthenticate(messageData)
	case messages.MessageTypeGoodbye:
		return parsers.ProtobufToGoodbye(messageData)
	case messages.MessageTypeCall:
		return parsers.ProtobufToCall(messageData, payloadData)
	case messages.MessageTypeInvocation:
		return parsers.ProtobufToInvocation(messageData, payloadData)
	case messages.MessageTypeYield:
		return parsers.ProtobufToYield(messageData, payloadData)
	case messages.MessageTypeResult:
		return parsers.ProtobufToResult(messageData, payloadData)
	case messages.MessageTypeRegister:
		return parsers.ProtobufToRegister(messageData)
	case messages.MessageTypeRegistered:
		return parsers.ProtobufToRegistered(messageData)
	case messages.MessageTypeUnregister:
		return parsers.ProtobufToUnregister(messageData)
	case messages.MessageTypeUnregistered:
		return parsers.ProtobufToUnregistered(messageData)
	case messages.MessageTypeAbort:
		return parsers.ProtobufToAbort(messageData, payloadData)
	case messages.MessageTypeCancel:
		return parsers.ProtobufToCancel(messageData)
	case messages.MessageTypeInterrupt:
		return parsers.ProtobufToInterrupt(messageData)
	case messages.MessageTypePublish:
		return parsers.ProtobufToPublish(messageData, payloadData)
	case messages.MessageTypePublished:
		return parsers.ProtobufToPublished(messageData)
	case messages.MessageTypeEvent:
		return parsers.ProtobufToEvent(messageData, payloadData)
	case messages.MessageTypeSubscribe:
		return parsers.ProtobufToSubscribe(messageData)
	case messages.MessageTypeSubscribed:
		return parsers.ProtobufToSubscribed(messageData)
	case messages.MessageTypeUnsubscribe:
		return parsers.ProtobufToUnsubscribe(messageData)
	case messages.MessageTypeUnsubscribed:
		return parsers.ProtobufToUnsubscribed(messageData)
	case messages.MessageTypeError:
		return parsers.ProtobufToError(messageData, payloadData)
	default:
		return nil, fmt.Errorf("unknown message type: %v", bytes[0])
	}
}

func (p *ProtobufSerializer) Static() bool {
	return true
}
