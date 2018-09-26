package queueing

import (
	"framework/message"
	"fmt"
	"framework/components/queueing/queueingclientproxy"
	"shared/parameters"
)

type QueueingService struct{}

var Queues = map[string]chan string{}

func (q QueueingService) I_PosInvP(msg *message.Message) {
	op := msg.Payload.(message.Invocation).Op
	args := msg.Payload.(message.Invocation).Args

	switch op {
	case "publish":
		argsX := args.([]interface{})
		topic := argsX[0].(string)
		m := argsX[1].(string)
		r := q.Publish(topic, m)
		msgRep := message.Message{Payload: r}
		*msg = msgRep
	case "consume":
		argsX := args.([]interface{})
		topic := argsX[0].(string)
		r := q.Consume(topic)
		msgRep := message.Message{Payload: r}
		*msg = msgRep
	default:
		fmt.Println("Queueing:: Operation " + op + " is not implemented by Queue Server")
	}
}

func (QueueingService) Publish(topic string, msg string) bool {
	r := false

	if _, ok := Queues[topic]; !ok {
		Queues[topic] = make(chan string, parameters.QUEUE_SIZE)
	}

	if len(Queues[topic]) < parameters.QUEUE_SIZE {
		Queues[topic] <- msg
		r = true
	} else {
		r = false
	}
	return r
}

func (QueueingService) Consume(topic string) string {
	r := ""
	if _, ok := Queues[topic]; !ok {
		Queues[topic] = make(chan string, parameters.QUEUE_SIZE)
	}
	if len(Queues[topic]) == 0 {
		r = "EMPTY QUEUE"
	} else {
		r = <-Queues[topic]
	}
	return r
}

func LocateQueueing() queueingclientproxy.QueueingClientProxy {
	p := queueingclientproxy.QueueingClientProxy{Host: parameters.QUEUEING_HOST, Port: parameters.QUEUEING_PORT}
	return p
}
