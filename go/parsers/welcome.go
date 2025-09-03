package parsers

import (
	"google.golang.org/protobuf/proto"

	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-protobuf/go/gen"
)

type welcome struct {
	gen *gen.Welcome
}

func NewWelcomeFields(call *gen.Welcome) messages.WelcomeFields {
	return &welcome{
		gen: call,
	}
}

func (w *welcome) SessionID() uint64 {
	return w.gen.GetSessionId()
}

func (w *welcome) Details() map[string]any {
	details := map[string]any{
		"authid":       w.gen.Authid,
		"authrole":     w.gen.Authrole,
		"authmethod":   w.gen.Authmethod,
		"authprovider": w.gen.Authprovider,
	}

	roles := map[string]any{}

	genRoles := w.gen.GetRoles()
	if genRoles == nil {
		return details
	}

	if r := genRoles.GetDealer(); r != nil {
		dealer := map[string]any{}
		if r.GetCallerIdentification() {
			dealer["caller_identification"] = true
		}
		if r.GetCallTimeout() {
			dealer["call_timeout"] = true
		}
		if r.GetCallCanceling() {
			dealer["call_canceling"] = true
		}
		if r.GetProgressiveCallResults() {
			dealer["progressive_call_results"] = true
		}
		if r.GetPatternBasedRegistration() {
			dealer["pattern_based_registration"] = true
		}
		if r.GetSharedRegistration() {
			dealer["shared_registration"] = true
		}
		roles["dealer"] = dealer
	}

	if r := genRoles.GetBroker(); r != nil {
		broker := map[string]any{}
		if r.GetPublisherIdentification() {
			broker["publisher_identification"] = true
		}
		if r.GetPublisherExclusion() {
			broker["publisher_exclusion"] = true
		}
		if r.GetAcknowledgeEventReceived() {
			broker["acknowledge_event_received"] = true
		}
		if r.GetPatternBasedSubscription() {
			broker["pattern_based_subscription"] = true
		}
		roles["broker"] = broker
	}

	details["roles"] = roles

	return details
}

func WelcomeToProtobuf(welcome *messages.Welcome) ([]byte, error) {
	msg := &gen.Welcome{
		SessionId: welcome.SessionID(),
	}

	if authid, ok := welcome.Details()["authid"].(string); ok {
		msg.Authid = authid
	}

	if authrole, ok := welcome.Details()["authrole"].(string); ok {
		msg.Authrole = authrole
	}

	if authmethod, ok := welcome.Details()["authmethod"].(string); ok {
		msg.Authmethod = authmethod
	}

	if authprovider, ok := welcome.Details()["authprovider"].(string); ok {
		msg.Authprovider = authprovider
	}

	roles := &gen.Welcome_Roles{}
	if rolesMap, ok := welcome.Details()["roles"].(map[string]any); ok {
		if dealerMap, ok := rolesMap["dealer"].(map[string]any); ok {
			dealer := &gen.Welcome_Roles_Dealer{}
			if v, ok := dealerMap["caller_identification"].(bool); ok {
				dealer.CallerIdentification = v
			}
			if v, ok := dealerMap["call_timeout"].(bool); ok {
				dealer.CallTimeout = v
			}
			if v, ok := dealerMap["call_canceling"].(bool); ok {
				dealer.CallCanceling = v
			}
			if v, ok := dealerMap["progressive_call_results"].(bool); ok {
				dealer.ProgressiveCallResults = v
			}
			if v, ok := dealerMap["pattern_based_registration"].(bool); ok {
				dealer.PatternBasedRegistration = v
			}
			roles.Dealer = dealer
		}

		if brokerMap, ok := rolesMap["broker"].(map[string]any); ok {
			broker := &gen.Welcome_Roles_Broker{}
			if v, ok := brokerMap["publisher_identification"].(bool); ok {
				broker.PublisherIdentification = v
			}
			if v, ok := brokerMap["publisher_exclusion"].(bool); ok {
				broker.PublisherExclusion = v
			}
			if v, ok := brokerMap["acknowledge_event_received"].(bool); ok {
				broker.AcknowledgeEventReceived = v
			}
			if v, ok := brokerMap["pattern_based_subscription"].(bool); ok {
				broker.PatternBasedSubscription = v
			}
			roles.Broker = broker
		}
	}

	msg.Roles = roles

	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return PrependHeader(messages.MessageTypeWelcome, data, nil), nil
}

func ProtobufToWelcome(data []byte) (*messages.Welcome, error) {
	msg := &gen.Welcome{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}

	return messages.NewWelcomeWithFields(NewWelcomeFields(msg)), nil
}
