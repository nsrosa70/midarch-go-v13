package naming

import (
	"framework/message"
	"transport/myRPC/ior"
	"framework/components/naming/namingimpl"
)

type NamingService struct{}
var ns = namingimpl.NamingImpl{}

func (n NamingService) I_PosInvP(msg *message.Message) {
	op := msg.Payload.(message.Invocation).Op

	switch op {
	case "register":
		args := msg.Payload.(message.Invocation).Args
	    argsX := args.([]interface{})
	    _p1 := argsX[0].(string)
		_p2 := argsX[1].(map[string]interface{})
		ior := ior.IOR{Host:_p2["Host"].(string),Port:int(_p2["Port"].(float64)),Id:int(_p2["Id"].(float64)),Proxy:_p2["Proxy"].(string)}
		_r := ns.Register(_p1,ior)
		msgRep := message.Message{_r}
		*msg = msgRep
	case "lookup":
		args := msg.Payload.(message.Invocation).Args
		argsX := args.([]interface{})
		_p1 := argsX[0].(string)
		_r := n.Lookup(_p1)
		msgRep := message.Message{_r}
		*msg = msgRep
	case "list":
		_r := n.List()
		msgRep := message.Message{_r}
		*msg = msgRep
	}
}