package queueingclientproxy

import (
	"framework/message"
	"reflect"
	"fmt"
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
	msg :=reflect.ValueOf(args[0]).String()
	argsTemp := []interface{}{msg}

	requestHeader := message.RequestHeader{Operation:"publish"}
	requestBody := message.RequestBody{Args:argsTemp}
	miopHeader := message.MIOPHeader{Magic:"MIOP"}
	miopBody := message.MIOPBody{RequestHeader:requestHeader,RequestBody:requestBody}

	miop := message.MIOP{Header:miopHeader,Body:miopBody}
	inv := message.ToCRH{Host:n.Host,Port:n.Port,MIOP:miop}

	reqMsg := message.Message{inv}

	chIn <- reqMsg

	fmt.Println("QueueingClientProxy:: HERE")
	repMsg := <-chOut
	payload := repMsg.Payload.(map[string]interface{})
	reply := payload["Reply"].(bool)
	return reply
}

func (n QueueingClientProxy) Subscribe()[]interface{} { //TODO
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