package components

import (
	"framework/messages"
	"shared/shared"
	"fmt"
	"os"
	"shared/parameters"
)

type NotificationConsumer struct{}

func (NC NotificationConsumer) I_PosInvP(msg *messages.SAMessage, r *bool) {
	inv := msg.Payload.(messages.Invocation)

	switch inv.Op {
	case "Publish":
		_args := inv.Args.([]interface{})
		_topic := _args[0].(string)
		_msg := _args[1].(map[string]interface{})
		_msgHeader := _msg["Header"].(map[string]interface{})
		_headerDestination := _msgHeader["Destination"].(string)
		_msgPayload := _msg["PayLoad"].(string)
		_r := NC.Publish(_topic, messages.MessageMOM{Header: messages.Header{Destination: _headerDestination}, PayLoad: _msgPayload})

		_ter := shared.QueueingTermination{_r}
		*msg = messages.SAMessage{_ter}
	case "Consume":
		_args := inv.Args.([]interface{})
		_topic := _args[0].(string)
		_r := NC.Consume(_topic)

		_ter := shared.QueueingTermination{_r}
		*msg = messages.SAMessage{_ter}

	default:
		fmt.Println("NotificationEngine:: Operation " + inv.Op + " is not implemented by NotificationEngine")
		os.Exit(0)
	}
}

func (NotificationConsumer) Publish(topic string, msg messages.MessageMOM) bool {
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

func (NotificationConsumer) Consume(topic string) messages.MessageMOM {
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
