package components

import (
	"fmt"
	"os"
	"framework/messages"
	"shared/shared"
	"shared/parameters"
)

type NotificationEngine struct{}

var SubscribersNE = map[string][]SubscriberRecord{}
var Topics = map[string]chan messages.MessageMOM{}
var MsgToBeNotified string
var TopicToBePublished string

func (NotificationEngine) I_PosInvP(msg *messages.SAMessage, r *bool) {

	switch msg.Payload.(messages.Invocation).Op {
	case "Subscribe":
	case "Unsubscribe":
	case "Publish":
	case "GetSubscribers":
	default:
		fmt.Println("NotificationEngine:: Operation " + msg.Payload.(messages.Invocation).Op + " is not implemented by NotificationEngine")
		os.Exit(0)
	}
}

func (NotificationEngine) I_Subscribe(msg *messages.SAMessage, r *bool) { // Subscribe Manager
	inv := msg.Payload.(messages.Invocation)

	switch inv.Op {
	case "Subscribe":
		*r = true
	default:
		*r = false
	}
}

func (NotificationEngine) I_Unsubscribe(msg *messages.SAMessage, r *bool) { // Subscribe Manager
	inv := msg.Payload.(messages.Invocation)

	switch inv.Op {
	case "Unsubscribe":
		*r = true
	default:
		*r = false
	}
}

func (NE NotificationEngine) I_Publish(msg *messages.SAMessage, r *bool) { // NOTIFICATION CONSUMER
	inv := msg.Payload.(messages.Invocation)

	switch inv.Op {
	case "Publish":
		*r = true
	default:
		*r = false
	}
}

func (NotificationEngine) I_GetSubscribers(msg *messages.SAMessage, r *bool) { // Subscribe Manager
	inv := msg.Payload.(messages.Invocation)

	switch inv.Op {
	case "GetSubscribers":
		*r = true
	default:
		*r = false
	}
}

func (NotificationEngine) I_GetSubs(msg *messages.SAMessage, r *bool) { // Subscribe Manager
	inv := messages.Invocation{Op: "GetSubscribers"}

	*msg = messages.SAMessage{Payload: inv}
}

func (NotificationEngine) I_GetResSubs(msg *messages.SAMessage, r *bool) { // Subscribe Manager
	ter := msg.Payload.(shared.QueueingTermination)

	SubscribersNE = ter.R.(map[string][]SubscriberRecord)
}

func (NE NotificationEngine) I_Pub(msg *messages.SAMessage, r *bool) {
	inv := msg.Payload.(messages.Invocation)

	_args := inv.Args.([]interface{})
	_topic := _args[0].(string)
	TopicToBePublished = _topic
	_msg := _args[1].(map[string]interface{})
	_msgHeader := _msg["Header"].(map[string]interface{})
	_headerDestination := _msgHeader["Destination"].(string)
	_msgPayload := _msg["PayLoad"].(string)
	MsgToBeNotified = _msgPayload
	_msgPub := messages.MessageMOM{Header: messages.Header{Destination: _headerDestination}, PayLoad: _msgPayload}

	_r := NE.Publish(_topic, _msgPub)

	_ter := shared.QueueingTermination{_r}
	*msg = messages.SAMessage{_ter}
}

func (NotificationEngine) I_Notify(msg *messages.SAMessage, r *bool) {

	tempSubs := filterSubscribers(TopicToBePublished)
	args := []interface{}{MsgToBeNotified, tempSubs}
	inv := messages.Invocation{Op: "Notify", Args: args}
	*msg = messages.SAMessage{Payload: inv}
}

func (NotificationEngine) I_Consume(msg *messages.SAMessage, r *bool) { // TODO
	*r = false
}

func (NotificationEngine) Consume(topic string) messages.MessageMOM {
	r := messages.MessageMOM{}
	if _, ok := Topics[topic]; !ok {
		Topics[topic] = make(chan messages.MessageMOM, parameters.QUEUE_SIZE)
	}
	if len(Topics[topic]) == 0 {
		r = messages.MessageMOM{Header: messages.Header{Destination: topic}, PayLoad: "QUEUE EMPTY"} // TODO
	} else {
		r = <-Topics[topic]
	}
	return r
}

func (NotificationEngine) Publish(topic string, msg messages.MessageMOM) bool {
	r := false

	if _, ok := Topics[topic]; !ok {
		Topics[topic] = make(chan messages.MessageMOM, parameters.QUEUE_SIZE)
	}

	if len(Topics[topic]) < parameters.QUEUE_SIZE {
		Topics[topic] <- msg
		r = true
	} else {
		r = false
	}
	return r
}

func filterSubscribers(topic string) []SubscriberRecord {
	r := []SubscriberRecord{}
	for i := range SubscribersNE {
		if i == topic {
			r = SubscribersNE[i]
		}
	}
	return r
}
