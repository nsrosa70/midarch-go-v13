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
	inv := msg.Payload.(shared.QueueingInvocation)

	switch inv.Op {
	case "Publish":
		_topic := inv.Args[0].(string)
		_msg := inv.Args[1].(MessageMOM)
		_r := QS.Publish(_topic,_msg)

		_ter := shared.QueueingTermination{_r}
		*msg = message.Message{_ter}
	default:
		fmt.Println("NamingInvoker:: Operation " + inv.Op + " is not implemented by Naming Server")
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