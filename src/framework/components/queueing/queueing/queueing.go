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

	// recover parameters
	op := msg.Payload.(message.Invocation).Op
	args := msg.Payload.(message.Invocation).Args

	switch op {
	case "publish":
		argsX := args.([]interface{})
		fmt.Println(argsX)
		topic := argsX[0].(string)
		m := argsX[1].(string)
		r := q.Publish(topic,m)
		msgRep := message.Message{r}
		*msg = msgRep
	default:
		fmt.Println("Queueing:: Operation "+op+" is not implemented by Queue Server")
	}
}

func (QueueingService) Publish(topic string, msg string) int {

	if _, ok := Queues[topic]; !ok {
		Queues[topic] = make(chan string, parameters.QUEUE_SIZE)
	}

	Queues[topic] <- msg

	fmt.Println("Queueing:: Publish :: "+topic+" "+msg)

	return len(Queues[topic])
}

func LocateQueueing(host string, port int) queueingclientproxy.QueueingClientProxy{
	p := queueingclientproxy.QueueingClientProxy{Host: host, Port: port}
	return p
}
