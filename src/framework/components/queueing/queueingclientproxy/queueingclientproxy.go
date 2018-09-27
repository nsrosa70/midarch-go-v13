package queueingclientproxy

import (
	"framework/message"
)

type QueueingClientProxy struct {
	Host string
	Port int
}

var i_PreInvR = make(chan message.Message)
var i_PosTerR = make(chan message.Message)

func (n QueueingClientProxy) Publish(_p1 string, _p2 string) bool {
	_args := []interface{}{_p1, _p2}
	_reqMsg := message.Message{message.Invocation{Host: n.Host, Port: n.Port, Op: "Publish", Args: _args}}

	i_PreInvR <- _reqMsg
	_repMsg := <-i_PosTerR

	_payload := _repMsg.Payload.(map[string]interface{})
	_reply := _payload["Reply"].(bool)
	return _reply
}

func (n QueueingClientProxy) Consume(_p1 string) string {
	_args := []interface{}{_p1}
	_reqMsg := message.Message{message.Invocation{Host: n.Host, Port: n.Port, Op: "Consume", Args: _args}}

	i_PreInvR <- _reqMsg
	_repMsg := <-i_PosTerR

	_payload := _repMsg.Payload.(map[string]interface{})
	_reply := _payload["Reply"].(string)
	return _reply
}

func (QueueingClientProxy) I_PreInvR(msg *message.Message) {
	*msg = <-i_PreInvR
}

func (QueueingClientProxy) I_PosTerR(msg *message.Message) {
	i_PosTerR <- *msg
}
