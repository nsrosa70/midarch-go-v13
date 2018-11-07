package components

import (
	"fmt"
	"os"
	"framework/messages"
	"shared/shared"
	"shared/parameters"
)

type NotificationEngine struct{}

var Subs = map[string][]string{}
var Queues = map[string]chan messages.MessageMOM{}
var MsgToBeNotified string
var TopicToBePublished string

func (NotificationEngine) I_PosInvP(msg *messages.SAMessage, r *bool) {

	switch msg.Payload.(messages.Invocation).Op {
	case "Subscribe":
	case "Unsubscribe":
	case "Publish":
	case "Consume":
	case "GetSubscribers":
	default:
		fmt.Println("NotificationEngine:: Operation " + msg.Payload.(messages.Invocation).Op + " is not implemented by NotificationEngine")
		os.Exit(0)
	}
}

func (NotificationEngine) I_GetSubs(msg *messages.SAMessage, r *bool) { // Subscribe Manager
	inv := messages.Invocation{Op: "GetSubscribers"}

	*msg = messages.SAMessage{Payload: inv}
}

func (NotificationEngine) I_GetResSubs(msg *messages.SAMessage, r *bool) { // Subscribe Manager
	ter := msg.Payload.(shared.QueueingTermination)

	Subs = ter.R.(map[string][]string)
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
		_args := inv.Args.([]interface{})
		_topic := _args[0].(string)
		_msg := _args[1].(map[string]interface{})
		_msgHeader := _msg["Header"].(map[string]interface{})
		_headerDestination := _msgHeader["Destination"].(string)
		_msgPayload := _msg["PayLoad"].(string)
		MsgToBeNotified = _msgPayload
		_msgPub := messages.MessageMOM{Header: messages.Header{Destination: _headerDestination}, PayLoad: _msgPayload}
		_r := NE.Publish(_topic, _msgPub)

		_ter := shared.QueueingTermination{_r}
		*msg = messages.SAMessage{_ter}
		*r = true
	case "Consume":
		_args := inv.Args.([]interface{})
		_topic := _args[0].(string)
		_r := NE.Consume(_topic)

		_ter := shared.QueueingTermination{_r}
		*msg = messages.SAMessage{_ter}
		*r = true

	default:
		*r = false
	}
}

func (NotificationEngine) I_Notify(msg *messages.SAMessage, r *bool){
    tempSubs := filterSubscribers(TopicToBePublished)
    args := []interface{}{MsgToBeNotified,tempSubs}
	inv := messages.Invocation{Op: "Notify",Args:args}
	*msg = messages.SAMessage{Payload: inv}
}

func filterSubscribers(topic string) []string{
	tempSubs := []string{}
	for i:= range Subs{
		if i == topic{
			tempSubs = Subs[i]
		}
	}
	return tempSubs
}

func (NotificationEngine) Consume(topic string) messages.MessageMOM {
	r := messages.MessageMOM{}
	if _, ok := Queues[topic]; !ok {
		Queues[topic] = make(chan messages.MessageMOM, parameters.QUEUE_SIZE)
	}
	if len(Queues[topic]) == 0 {
		r = messages.MessageMOM{Header: messages.Header{Destination: topic}, PayLoad: "QUEUE EMPTY"} // TODO
	} else {
		r = <-Queues[topic]
	}
	return r
}

func (NotificationEngine) I_Consume(msg *messages.SAMessage, r *bool) { // NOTIFICATION CONSUMER
	inv := msg.Payload.(messages.Invocation)

	// TODO
	switch inv.Op {
	case "Consume":
		*r = true
	default:
		*r = false
	}
}

func (NotificationEngine) Publish(topic string, msg messages.MessageMOM) bool {
	r := false

	if _, ok := Queues[topic]; !ok {
		Queues[topic] = make(chan messages.MessageMOM, parameters.QUEUE_SIZE)
	}

	if len(Queues[topic]) < parameters.QUEUE_SIZE {
		Queues[topic] <- msg
		r = true
	} else {
		r = false
	}

	return r
}
