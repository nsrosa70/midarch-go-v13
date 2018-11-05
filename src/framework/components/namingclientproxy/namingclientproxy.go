package namingclientproxy

import (
	"framework/messages"
	"reflect"
	"framework/element"
	"framework/proxy"
)

type NamingClientProxy struct {
	Host string
	Port int
}

var i_PreInvRNamingClientProxy = make(chan messages.SAMessage)
var i_PosTerRNamingClientProxy = make(chan messages.SAMessage)

func (n NamingClientProxy) Register(_p1 string, _p2 interface{}) bool {
	_p3 := reflect.ValueOf(_p2).FieldByName("Host").String()
	_p4 := int(reflect.ValueOf(_p2).FieldByName("Port").Int())
	_p5 := reflect.TypeOf(_p2).String()
	_args := []interface{}{_p1, element.IOR{Host: _p3, Port: _p4, Proxy: _p5, Id: 1313}}
	_reqMsg := messages.SAMessage{messages.Invocation{Host: n.Host, Port: n.Port, Op: "Register", Args: _args}}
	i_PreInvRNamingClientProxy <- _reqMsg

	_repMsg := <-i_PosTerRNamingClientProxy
	_payload := _repMsg.Payload.(map[string]interface{})
	_reply := _payload["Reply"].(bool)
	return _reply
}

func (n NamingClientProxy) List() []interface{} {
	_reqMsg := messages.SAMessage{messages.Invocation{Host: n.Host, Port: n.Port, Op: "List"}}
	i_PreInvRNamingClientProxy <- _reqMsg

	_repMsg := <-i_PosTerRNamingClientProxy
	_payload := _repMsg.Payload.(map[string]interface{})
	_reply := _payload["Reply"].([]interface{})
	return _reply
}

func (n NamingClientProxy) Lookup(_p1 string) interface{} {
	_args := []string{_p1}
	_reqMsg := messages.SAMessage{messages.Invocation{Host: n.Host, Port: n.Port, Op: "Lookup", Args: _args}}
	i_PreInvRNamingClientProxy <- _reqMsg

	_repMsg := <-i_PosTerRNamingClientProxy
	_payload := _repMsg.Payload.(map[string]interface{})
	_reply := _payload["Reply"].(map[string]interface{})
	_proxyName := _reply["Proxy"].(string)
	_port := int64(_reply["Port"].(float64))
	_host := _reply["Host"].(string)
	//_p := reflect.TypeOf(libraries.Repository[_proxyName].Go)
	_p := proxy.ProxyLibrary[_proxyName]

	proxyPointer := reflect.New(_p)
	proxyValue := proxyPointer.Elem()
	proxyValue.FieldByName("Host").SetString(_host)
	proxyValue.FieldByName("Port").SetInt(_port)
	proxyInterface := proxyValue.Interface()

	return proxyInterface
}

func (NamingClientProxy) I_PreInvR(msg *messages.SAMessage,r *bool) {
	*msg = <-i_PreInvRNamingClientProxy
}

func (NamingClientProxy) I_PosTerR(msg *messages.SAMessage,r *bool) {
	i_PosTerRNamingClientProxy <- *msg
}
