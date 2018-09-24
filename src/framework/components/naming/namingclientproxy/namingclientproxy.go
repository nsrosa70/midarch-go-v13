package namingclientproxy

import (
	"framework/message"
	"framework/components/clientproxy/proxy"
	"reflect"
	"transport/middleware/ior"
)

type NamingClientProxy struct {
	Host string
	Port int
}

//var reqMsg message.Message
//var repMsg message.Message
//var opRequested = false
//var opFinished = false

var chIn = make(chan message.Message)
var chOut = make(chan message.Message)

func (n NamingClientProxy) Register(args ... interface{}) bool {

	port := int(reflect.ValueOf(args[1]).FieldByName("Port").Int())
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

func (n NamingClientProxy) List() []interface{} {
	inv := message.Invocation{Host: n.Host, Port: n.Port, Op: "list"}
	reqMsg := message.Message{inv}

	chIn <- reqMsg
	repMsg := <-chOut
	payload := repMsg.Payload.(map[string]interface{})
	reply := payload["Reply"].([]interface{})
	return reply
}

func (n NamingClientProxy) Lookup(args ... interface{}) interface{} {

	inv := message.Invocation{Host: n.Host, Port: n.Port, Op: "lookup", Args: args}
	reqMsg := message.Message{inv}

	chIn <- reqMsg
	repMsg := <-chOut
			payload   := repMsg.Payload.(map[string]interface{})
			reply     := payload["Reply"].(map[string]interface{})
			proxyName := reply["Proxy"].(string)
			port      := int64(reply["Port"].(float64))
			host      := reply["Host"].(string)
			p         := proxy.ProxyLibrary[proxyName]

			proxyPointer := reflect.New(p)
			proxyValue := proxyPointer.Elem()
			proxyValue.FieldByName("Host").SetString(host)
			proxyValue.FieldByName("Port").SetInt(port)
			proxyInterface := proxyValue.Interface()

			return proxyInterface
}

func (NamingClientProxy) I_PreInvR(msg *message.Message) {
	*msg = <-chIn
}

func (NamingClientProxy) I_PosTerR(msg *message.Message) {
	chOut <- *msg
}