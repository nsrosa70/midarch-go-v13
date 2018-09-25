package queueing

import (
	"framework/message"
	"fmt"
	"framework/components/queueing/queueingclientproxy"
)

type QueueingService struct{}

//var Repo = map[string]ior.IOR{}

func (q QueueingService) I_PosInvP(msg *message.Message) {

	// recover parameters
	op := msg.Payload.(message.Invocation).Op
	args := msg.Payload.(message.Invocation).Args

	switch op {
	case "publish":
		fmt.Println("Queueing:: Publish")
		argsX := args.([]interface{})
		fmt.Println(argsX)
		p1 := argsX[0].(string)
		r := q.Publish(p1)
		msgRep := message.Message{r}
		*msg = msgRep
	default:
		fmt.Println("Queueing:: Operation "+op+" is not implemented by Queue Server")
	}
}

func (QueueingService) Publish(msg string) bool {

	return true // TODO
}

func LocateQueueing(host string, port int) queueingclientproxy.QueueingClientProxy{
	p := queueingclientproxy.QueueingClientProxy{Host: host, Port: port}
	return p
}
