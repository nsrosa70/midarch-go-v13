package queueingclientproxy

import (
	"framework/message"
	"reflect"
	"transport/myRPC/ior"
)

type QueueingClientProxy struct {
	Host string
	Port int
}

//var reqMsg message.Message
//var repMsg message.Message
//var opRequested = false
//var opFinished = false

var chIn = make(chan message.Message)
var chOut = make(chan message.Message)

func (n QueueingClientProxy) Publish(args ... interface{}) bool {
	message := int(reflect.ValueOf(args[1]).FieldByName("Port").Int())
	host := reflect.ValueOf(args[1]).FieldByName("Host").String()
	proxy := reflect.TypeOf(args[1]).String()
	ior := ior.IOR{Host: host, Port: port, Proxy: proxy, Id: 1313} // TODO
	argsTemp := []interface{}{args[0], ior}
	inv := message.Invocation{Host: n.Host, Port: n.Port, Op: "register", Args: argsTemp}
	reqMsg := message.Message{inv}

	chIn <- reqMsg
	repMsg := <-chOut
	payload := repMsg.Payload.(map[string]interface{})
	reply := payload["Reply"].(bool)
	return reply
}

func (n QueueingClientProxy) Subscribe() []interface{} {
	inv := message.Invocation{Host: n.Host, Port: n.Port, Op: "list"}
	reqMsg := message.Message{inv}

	chIn <- reqMsg
	repMsg := <-chOut
	payload := repMsg.Payload.(map[string]interface{})
	reply := payload["Reply"].([]interface{})
	return reply
}

func (QueueingClientProxy) I_PreInvR(msg *message.Message) {
	*msg = <-chIn
}

func (QueueingClientProxy) I_PosTerR(msg *message.Message) {
	chOut <- *msg
}