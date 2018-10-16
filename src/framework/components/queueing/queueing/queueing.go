package queueing

import (
	"shared/parameters"
)

type Queueing struct {}

type MessageMOM struct{
	Header string
	PayLoad string
}

var Queues = map[string]chan MessageMOM{}

func (Queueing) Publish(topic string, msg MessageMOM) bool {
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

func (Queueing) Consume(topic string) MessageMOM {
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
