package queueing

import (
	"shared/parameters"
)

type Queueing1 struct {}

type MessageMOM1 struct{
	Header string
	PayLoad string
}

var Queues1 = map[string]chan MessageMOM{}

func (Queueing1) Publish(topic string, msg MessageMOM) bool {
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

func (Queueing1) Consume(topic string) MessageMOM {
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
