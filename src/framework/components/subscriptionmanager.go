package components

import (
	"framework/messages"
	"shared/shared"
	"fmt"
	"os"
)

type SubscriptionManager struct{}

var Subscribers = map[string][]string{} // TODO

func (c SubscriptionManager) I_PosInvP(msg *messages.SAMessage, r *bool) {
	inv := msg.Payload.(messages.Invocation)

	switch inv.Op {
	case "Subscribe":
		_args := inv.Args.([]interface{})
		_topic := _args[0].(string)
		_ip := _args[1].(string)
		_r := c.Subscribe(_topic, _ip)

		_ter := shared.QueueingTermination{_r}
		*msg = messages.SAMessage{_ter}
	case "Unsubscribe":
	case "GetSubscribers":
		_r := c.GetSubscribers()

		_ter := shared.QueueingTermination{_r}
		*msg = messages.SAMessage{_ter}

	default:
		fmt.Println("SubscriptionManager:: Operation " + inv.Op + " is not implemented by SubscriptionManager")
		os.Exit(0)
	}
}

func (SubscriptionManager) GetSubscribers() map[string][]string {

	return Subscribers
}

func (SubscriptionManager) Subscribe(topic string, ip string) bool {
	r := true

	if _, ok := Subscribers[topic]; !ok {
		Subscribers[topic] = []string{}
	}

	Subscribers[topic] = append(Subscribers[topic], ip)

	return r
}
