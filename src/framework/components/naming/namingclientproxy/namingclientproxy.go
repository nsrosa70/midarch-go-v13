package namingclientproxy

import (
	"framework/message"
	"framework/components/clientproxy/proxy"
	"reflect"
	"framework/components/naming/ior"
)

type NamingClientProxy struct {
	Host string
	Port int
}

var i_PreInvR = make(chan message.Message)
var i_PosTerR = make(chan message.Message)

func (e NamingClientProxy) Loop(channels map[string] chan message.Message) {
	var msgTerR, msgPreInvR message.Message
	for {
		select {
		case msgPreInvR = <-channels["I_PreInvR_namingproxy"]:
			e.I_PreInvR(&msgPreInvR)
		case channels["InvR"] <- msgPreInvR:
		case msgTerR = <-channels["TerR"]:
		case <-channels["I_PosTerR_namingproxy"]:
			e.I_PosTerR(&msgTerR)
		}
	}
}

func (n NamingClientProxy) Register(_p1 string, _p2 interface{}) bool {
	_p3 := reflect.ValueOf(_p2).FieldByName("Host").String()
	_p4 := int(reflect.ValueOf(_p2).FieldByName("Port").Int())
	_p5 := reflect.TypeOf(_p2).String()
	_args := []interface{}{_p1, ior.IOR{Host: _p3, Port: _p4, Proxy: _p5, Id: 1313}}
	_reqMsg := message.Message{message.Invocation{Host: n.Host, Port: n.Port, Op: "Register", Args: _args}}
	i_PreInvR <- _reqMsg

	_repMsg := <-i_PosTerR
	_payload := _repMsg.Payload.(map[string]interface{})
	_reply := _payload["Reply"].(bool)
	return _reply
}

func (n NamingClientProxy) List() []interface{} {
	_reqMsg := message.Message{message.Invocation{Host: n.Host, Port: n.Port, Op: "List"}}
	i_PreInvR <- _reqMsg

	_repMsg := <-i_PosTerR
	_payload := _repMsg.Payload.(map[string]interface{})
	_reply := _payload["Reply"].([]interface{})
	return _reply
}

func (n NamingClientProxy) Lookup(_p1 string) interface{} {
	_args := []string{_p1}
	_reqMsg := message.Message{message.Invocation{Host: n.Host, Port: n.Port, Op: "Lookup", Args: _args}}
	i_PreInvR <- _reqMsg

	_repMsg := <-i_PosTerR
	_payload := _repMsg.Payload.(map[string]interface{})
	_reply := _payload["Reply"].(map[string]interface{})
	_proxyName := _reply["Proxy"].(string)
	_port := int64(_reply["Port"].(float64))
	_host := _reply["Host"].(string)
	p := proxy.ProxyLibrary[_proxyName]

	proxyPointer := reflect.New(p)
	proxyValue := proxyPointer.Elem()
	proxyValue.FieldByName("Host").SetString(_host)
	proxyValue.FieldByName("Port").SetInt(_port)
	proxyInterface := proxyValue.Interface()

	return proxyInterface
}

func (NamingClientProxy) I_PreInvR(msg *message.Message) {
	*msg = <-i_PreInvR
}

func (NamingClientProxy) I_PosTerR(msg *message.Message) {
	i_PosTerR <- *msg
}
