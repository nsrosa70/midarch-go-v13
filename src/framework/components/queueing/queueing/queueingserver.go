package queueing

import (
	"shared/parameters"
	"framework/message"
	"fmt"
	"shared/shared"
	"os"
)

type QueueingServer struct {}

type MessageMOM struct{
	Header string
	PayLoad string
}

var Queues = map[string]chan MessageMOM{}

func (QS QueueingServer) I_PosInvP(msg *message.Message){
	inv := msg.Payload.(message.Invocation)

	switch inv.Op {
	case "Publish":
		_args := inv.Args.([]interface{})
		_topic := _args[0].(string)
		_msg := _args[1].(map[string]interface{})
		_r := QS.Publish(_topic,MessageMOM{Header:_msg["Header"].(string),PayLoad:_msg["PayLoad"].(string)})

		_ter := shared.QueueingTermination{_r}
		*msg = message.Message{_ter}
	case "Consume":
		_args := inv.Args.([]interface{})
		_topic := _args[0].(string)
		_r := QS.Consume(_topic)

		_ter := shared.QueueingTermination{_r}
		*msg = message.Message{_ter}

	default:
		fmt.Println("QueueingInvoker:: Operation " + inv.Op + " is not implemented by Queueing Server")
		os.Exit(0)
	}
}


func (QueueingServer) Publish(topic string, msg MessageMOM) bool {
	r := false

	if _, ok := Queues[topic]; !ok {
		Queues[topic] = make(chan MessageMOM, parameters.QUEUE_SIZE)
	}

	if len(Queues[topic]) < parameters.QUEUE_SIZE {
		Queues[topic] <- msg
		r = true
	} else {
		r = false
	}
	return r
}

func (QueueingServer) Consume(topic string) MessageMOM {
	r := MessageMOM{}
	if _, ok := Queues[topic]; !ok {
		Queues[topic] = make(chan MessageMOM, parameters.QUEUE_SIZE)
	}
	if len(Queues[topic]) == 0 {
		r = MessageMOM{Header:"header",PayLoad:"QUEUE EMPTY"}
	} else {
		r = <-Queues[topic]
	}
	return r
}