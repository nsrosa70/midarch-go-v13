package components

import (
	"framework/messages"
	"shared/net"
	"shared/Handlers"
)

type NotificationEngineClientProxy struct {
	Host string
	Port int
}

var i_PreInvRNotificationEngineClientProxy = make(chan messages.SAMessage)
var i_PosTerRNotificationEngineClientProxy = make(chan messages.SAMessage)
var HandlersProxy = make(map[string]Handlers.HandlerNotify,10)  // TODO

func (p NotificationEngineClientProxy) Subscribe(_p1 string) (Handlers.HandlerNotify,bool) {
	_p2 := netshared.ResolveHostIp()             // host
	_p3 := netshared.NextPortTCPAvailable()      // port
	_args := []interface{}{_p1,_p2,_p3}
	_reqMsg := messages.SAMessage{messages.Invocation{Host: p.Host, Port: p.Port, Op: "Subscribe", Args: _args}}

	i_PreInvRNotificationEngineClientProxy <- _reqMsg
	_repMsg := <-i_PosTerRNotificationEngineClientProxy

	_payload := _repMsg.Payload.(map[string]interface{})
	_reply := _payload["Reply"].(map[string]interface{})
	_r := _reply["R"].(bool)

	// Include the new handler associated to the topic [one handler per topic]
	if _,ok := HandlersProxy[_p1]; !ok{
		HandlersProxy[_p1] = Handlers.HandlerNotify{Host:_p2,Port:_p3}
	}
	HandlersProxy[_p1].StartHandler()

	return HandlersProxy[_p1],_r
}

func (p NotificationEngineClientProxy) Unsubscribe(_p1 string) (bool) {
	_p2 := HandlersProxy[_p1].Host
	_p3 := HandlersProxy[_p1].Port
	_args := []interface{}{_p1,_p2,_p3}
	_reqMsg := messages.SAMessage{messages.Invocation{Host: p.Host, Port: p.Port, Op: "Unsubscribe", Args: _args}}

	i_PreInvRNotificationEngineClientProxy <- _reqMsg
	_repMsg := <-i_PosTerRNotificationEngineClientProxy

	_payload := _repMsg.Payload.(map[string]interface{})
	_reply := _payload["Reply"].(map[string]interface{})
	_r := _reply["R"].(bool)

	// Stop Handler // TODO

	return _r
}

func (n NotificationEngineClientProxy) Publish(_p1 string, _p2 interface{}) bool {
	_tempP2 := messages.MessageMOM{Header:messages.Header{"Header"},PayLoad:_p2.(string)}
	_args := []interface{}{_p1, _tempP2}
	_reqMsg := messages.SAMessage{messages.Invocation{Host: n.Host, Port: n.Port, Op: "Publish", Args: _args}}

	i_PreInvRNotificationEngineClientProxy <- _reqMsg
	_repMsg := <-i_PosTerRNotificationEngineClientProxy

	_payload := _repMsg.Payload.(map[string]interface{})
	_reply := _payload["Reply"].(map[string]interface{})
	_r := _reply["R"].(bool)

	return _r
	}

func (NE NotificationEngineClientProxy) Consume(_p1 string) messages.MessageMOM {
	_args := []interface{}{_p1}
	_reqMsg := messages.SAMessage{messages.Invocation{Host: NE.Host, Port: NE.Port, Op: "Consume", Args: _args}}

	i_PreInvRNotificationEngineClientProxy <- _reqMsg
	_repMsg := <-i_PosTerRNotificationEngineClientProxy

	_payload := _repMsg.Payload.(map[string]interface{})
	_reply := _payload["Reply"].(map[string]interface{})
	_rTemp := _reply["R"].(map[string]interface{})
	_msgHeader := _rTemp["Header"].(map[string]interface{})
	_headerDestination := _msgHeader["Destination"].(string)
	_msgPayload := _rTemp["PayLoad"].(string)
	_r := messages.MessageMOM{Header:messages.Header{Destination:_headerDestination},PayLoad:_msgPayload}

	return _r
}

func (NotificationEngineClientProxy) I_PreInvR(msg *messages.SAMessage,r *bool) {
	*msg = <-i_PreInvRNotificationEngineClientProxy
}

func (NotificationEngineClientProxy) I_PosTerR(msg *messages.SAMessage,r *bool) {
	i_PosTerRNotificationEngineClientProxy <- *msg
}
