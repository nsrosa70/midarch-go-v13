package queueingclientproxy

import (
	"framework/message"
	"framework/components/queueing/queueing"
)

type QueueingClientProxy struct {
	Host string
	Port int
}

var i_PreInvR = make(chan message.Message)
var i_PosTerR = make(chan message.Message)

func (n QueueingClientProxy) Publish(_p1 string, _p2 queueing.MessageMOM) bool {
	_args := []interface{}{_p1, _p2}
	_reqMsg := message.Message{message.Invocation{Host: n.Host, Port: n.Port, Op: "Publish", Args: _args}}

	i_PreInvR <- _reqMsg
	_repMsg := <-i_PosTerR

	_payload := _repMsg.Payload.(map[string]interface{})
	_reply := _payload["Reply"].(map[string]interface{})
	_r := _reply["R"].(bool)

	return _r
	}

func (n QueueingClientProxy) Consume(_p1 string) queueing.MessageMOM {
	_args := []interface{}{_p1}
	_reqMsg := message.Message{message.Invocation{Host: n.Host, Port: n.Port, Op: "Consume", Args: _args}}

	i_PreInvR <- _reqMsg
	_repMsg := <-i_PosTerR

	_payload := _repMsg.Payload.(map[string]interface{})
	_reply := _payload["Reply"].(map[string]interface{})
	_rTemp := _reply["R"].(map[string]interface{})
	_r := queueing.MessageMOM{Header:_rTemp["Header"].(string),PayLoad:_rTemp["PayLoad"].(string)}
	return _r
}

func (QueueingClientProxy) I_PreInvR(msg *message.Message) {
	*msg = <-i_PreInvR
}

func (QueueingClientProxy) I_PosTerR(msg *message.Message) {
	i_PosTerR <- *msg
}
