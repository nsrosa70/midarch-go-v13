package queueimpl

import (
	"shared/parameters"
)

type QueueServiceImpl struct {}
var Queues = map[string]chan string{}

func (QueueServiceImpl) Publish(topic string, msg string) bool {
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

func (QueueServiceImpl) Consume(topic string) string {
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
