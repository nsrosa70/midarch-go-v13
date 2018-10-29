package components

import (
	"framework/messages"
)

type NotificationEngineClientProxy struct {
	Host string
	Port int
}

var i_PreInvRNotificationEngineClientProxy = make(chan messages.SAMessage)
var i_PosTerRNotificationEngineClientProxy = make(chan messages.SAMessage)

func (p NotificationEngineClientProxy) Subscribe(_p1 string) bool {
	_args := []interface{}{_p1}
	_reqMsg := messages.SAMessage{messages.Invocation{Host: p.Host, Port: p.Port, Op: "Publish", Args: _args}}

	i_PreInvRNotificationEngineClientProxy <- _reqMsg
	_repMsg := <-i_PosTerRNotificationEngineClientProxy

	_payload := _repMsg.Payload.(map[string]interface{})
	_reply := _payload["Reply"].(map[string]interface{})
	_r := _reply["R"].(bool)

	return _r
}

func (n NotificationEngineClientProxy) Publish(_p1 string, _p2 messages.MessageMOM) bool {
	_args := []interface{}{_p1, _p2}
	_reqMsg := messages.SAMessage{messages.Invocation{Host: n.Host, Port: n.Port, Op: "Publish", Args: _args}}

	i_PreInvRNotificationEngineClientProxy <- _reqMsg
	_repMsg := <-i_PosTerRNotificationEngineClientProxy

	_payload := _repMsg.Payload.(map[string]interface{})
	_reply := _payload["Reply"].(map[string]interface{})
	_r := _reply["R"].(bool)

	return _r
	}

func (n NotificationEngineClientProxy) Consume(_p1 string) messages.MessageMOM {
	_args := []interface{}{_p1}
	_reqMsg := messages.SAMessage{messages.Invocation{Host: n.Host, Port: n.Port, Op: "Consume", Args: _args}}

	i_PreInvRNotificationEngineClientProxy <- _reqMsg
	_repMsg := <-i_PosTerRNotificationEngineClientProxy

	_payload := _repMsg.Payload.(map[string]interface{})
	_reply := _payload["Reply"].(map[string]interface{})
	_rTemp := _reply["R"].(map[string]interface{})
	_r := messages.MessageMOM{Header:_rTemp["Header"].(messages.Header),PayLoad:_rTemp["PayLoad"].(string)}
	return _r
}

func (NotificationEngineClientProxy) I_PreInvR(msg *messages.SAMessage) {
	*msg = <-i_PreInvRNotificationEngineClientProxy
}

func (NotificationEngineClientProxy) I_PosTerR(msg *messages.SAMessage) {
	i_PosTerRNotificationEngineClientProxy <- *msg
}
