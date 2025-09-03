package parsers

import (
	"google.golang.org/protobuf/proto"

	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-protobuf/go/gen"
)

type hello struct {
	gen *gen.Hello
}

func NewHelloFields(call *gen.Hello) messages.HelloFields {
	return &hello{
		gen: call,
	}
}

func (h *hello) Realm() string {
	return h.gen.Realm
}

func (h *hello) AuthID() string {
	return h.gen.GetAuthid()
}

func (h *hello) AuthMethods() []string {
	return h.gen.Authmethods
}

func (h *hello) AuthExtra() map[string]any {
	extra := make(map[string]any)
	if h.gen.PublicKey != "" {
		extra["pubkey"] = h.gen.PublicKey
	}
	return extra
}

func (h *hello) Roles() map[string]any {
	roles := map[string]any{}

	genRoles := h.gen.GetRoles()
	if genRoles == nil {
		return roles
	}

	if r := genRoles.GetCaller(); r != nil {
		roles["caller"] = map[string]any{
			"caller_identification":    r.GetCallerIdentification(),
			"call_timeout":             r.GetCallTimeout(),
			"call_canceling":           r.GetCallCanceling(),
			"progressive_call_results": r.GetProgressiveCallResults(),
		}
	}
	if r := genRoles.GetCallee(); r != nil {
		roles["callee"] = map[string]any{
			"caller_identification":      r.GetCallerIdentification(),
			"call_timeout":               r.GetCallTimeout(),
			"call_canceling":             r.GetCallCanceling(),
			"progressive_call_results":   r.GetProgressiveCallResults(),
			"pattern_based_registration": r.GetPatternBasedRegistration(),
			"shared_registration":        r.GetSharedRegistration(),
		}
	}
	if r := genRoles.GetPublisher(); r != nil {
		roles["publisher"] = map[string]any{
			"publisher_identification":   r.GetPublisherIdentification(),
			"publisher_exclusion":        r.GetPublisherExclusion(),
			"acknowledge_event_received": r.GetAcknowledgeEventReceived(),
		}
	}
	if r := genRoles.GetSubscriber(); r != nil {
		roles["subscriber"] = map[string]any{
			"publisher_identification":   r.GetPublisherIdentification(),
			"pattern_based_subscription": r.GetPatternBasedSubscription(),
		}
	}

	return roles

}

func HelloToProtobuf(h *messages.Hello) ([]byte, error) {
	msg := &gen.Hello{
		Realm:       h.Realm(),
		Authid:      h.AuthID(),
		Authmethods: h.AuthMethods(),
	}

	roles := &gen.Hello_Roles{}

	if callerMap, ok := h.Roles()["caller"].(map[string]any); ok {
		caller := &gen.Hello_Roles_Caller{}
		if v, ok := callerMap["caller_identification"].(bool); ok {
			caller.CallerIdentification = v
		}
		if v, ok := callerMap["call_timeout"].(bool); ok {
			caller.CallTimeout = v
		}
		if v, ok := callerMap["call_canceling"].(bool); ok {
			caller.CallCanceling = v
		}
		if v, ok := callerMap["progressive_call_results"].(bool); ok {
			caller.ProgressiveCallResults = v
		}
		roles.Caller = caller
	}

	if calleeMap, ok := h.Roles()["callee"].(map[string]any); ok {
		callee := &gen.Hello_Roles_Callee{}
		if v, ok := calleeMap["caller_identification"].(bool); ok {
			callee.CallerIdentification = v
		}
		if v, ok := calleeMap["call_timeout"].(bool); ok {
			callee.CallTimeout = v
		}
		if v, ok := calleeMap["call_canceling"].(bool); ok {
			callee.CallCanceling = v
		}
		if v, ok := calleeMap["progressive_call_results"].(bool); ok {
			callee.ProgressiveCallResults = v
		}
		if v, ok := calleeMap["pattern_based_registration"].(bool); ok {
			callee.PatternBasedRegistration = v
		}
		if v, ok := calleeMap["shared_registration"].(bool); ok {
			callee.SharedRegistration = v
		}
		roles.Callee = callee
	}

	if publisherMap, ok := h.Roles()["publisher"].(map[string]any); ok {
		publisher := &gen.Hello_Roles_Publisher{}
		if v, ok := publisherMap["publisher_identification"].(bool); ok {
			publisher.PublisherIdentification = v
		}
		if v, ok := publisherMap["publisher_exclusion"].(bool); ok {
			publisher.PublisherExclusion = v
		}
		if v, ok := publisherMap["acknowledge_event_received"].(bool); ok {
			publisher.AcknowledgeEventReceived = v
		}
		roles.Publisher = publisher
	}

	if subscriberMap, ok := h.Roles()["subscriber"].(map[string]any); ok {
		subscriber := &gen.Hello_Roles_Subscriber{}
		if v, ok := subscriberMap["publisher_identification"].(bool); ok {
			subscriber.PublisherIdentification = v
		}
		if v, ok := subscriberMap["pattern_based_subscription"].(bool); ok {
			subscriber.PatternBasedSubscription = v
		}
		roles.Subscriber = subscriber
	}

	msg.Roles = roles

	pubKey, ok := h.AuthExtra()["pubkey"].(string)
	if ok {
		msg.PublicKey = pubKey
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return PrependHeader(messages.MessageTypeHello, data, nil), nil
}

func ProtobufToHello(data []byte) (*messages.Hello, error) {
	msg := &gen.Hello{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}

	return messages.NewHelloWithFields(NewHelloFields(msg)), nil
}
