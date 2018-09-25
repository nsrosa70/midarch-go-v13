package queueingclientproxy

import (
	"framework/message"
	"reflect"
)

type QueueingClientProxy struct {
	Host string
	Port int
}

var chIn = make(chan message.Message)
var chOut = make(chan message.Message)

func (n QueueingClientProxy) Publish(args ... interface{}) bool {
	msg :=reflect.ValueOf(args[0]).String()
	argsTemp := []interface{}{msg}

	inv := message.Invocation{Host: n.Host, Port: n.Port, Op: "publish", Args: argsTemp}
	reqMsg := message.Message{inv}
	chIn <- reqMsg
	repMsg := <-chOut
	payload := repMsg.Payload.(map[string]interface{})
	reply := payload["Reply"].(bool)
	return reply
}

func (QueueingClientProxy) I_PreInvR(msg *message.Message) {
	*msg = <-chIn
}

func (QueueingClientProxy) I_PosTerR(msg *message.Message) {
	chOut <- *msg
}