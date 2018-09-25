package queueing

import (
	"framework/message"
	"fmt"
	"framework/components/queueing/queueingclientproxy"
)

type QueueingService struct{}

//var Repo = map[string]ior.IOR{}

func (QueueingService) I_PosInvP(msg *message.Message) {

	// recover parameters
	//op := msg.Payload.(message.Invocation).Op
	//args := msg.Payload.(message.Invocation).Args

	//switch op {
	//case "publish":
	//	fmt.Println("Naming:: Register")

	fmt.Println("Queueing:: HERE")
}

func LocateQueueing(host string, port int) queueingclientproxy.QueueingClientProxy{
	p := queueingclientproxy.QueueingClientProxy{Host: host, Port: port}
	return p
}
