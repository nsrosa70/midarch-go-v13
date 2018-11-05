package components

import (
	"framework/messages"
	"fmt"
	"os"
	"framework/element"
	"middlewareservices/naming"
)

type NamingInvoker struct{}

func (NamingInvoker) I_PosInvP(msg *messages.SAMessage, r *bool) {
	op := msg.Payload.(messages.MIOP).Body.RequestHeader.Operation
	switch op {
	case "Register":
		// Process request
		_args := msg.Payload.(messages.MIOP).Body.RequestBody.Args
		_argsX := _args.([]interface{})
		_p1 := _argsX[0].(string)
		_p2 := _argsX[1].(map[string]interface{})
		_ior := element.IOR{Host: _p2["Host"].(string), Port: int(_p2["Port"].(float64)), Id: int(_p2["Id"].(float64)), Proxy: _p2["Proxy"].(string)}
		_r := naming.NamingService{}.Register(_p1, _ior)

		// Send Reply
		_replyHeader := messages.ReplyHeader{Status: 1} // 1 - Success
		_replyBody := messages.ReplyBody{Reply: _r}
		_miopHeader := messages.MIOPHeader{Magic: "MIOP"}
		_miopBody := messages.MIOPBody{ReplyHeader: _replyHeader, ReplyBody: _replyBody}
		_miop := messages.MIOP{Header: _miopHeader, Body: _miopBody}
		*msg = messages.SAMessage{_miop}
	case "Lookup":
		// Process request
		_args := msg.Payload.(messages.MIOP).Body.RequestBody.Args
		_argsX := _args.([]interface{})
		_p1 := _argsX[0].(string)
		_r := naming.NamingService{}.Lookup(_p1)

		// Send Reply
		_replyHeader := messages.ReplyHeader{Status: 1} // 1 - Success
		_replyBody := messages.ReplyBody{Reply: _r}
		_miopHeader := messages.MIOPHeader{Magic: "MIOP"}
		_miopBody := messages.MIOPBody{ReplyHeader: _replyHeader, ReplyBody: _replyBody}
		_miop := messages.MIOP{Header: _miopHeader, Body: _miopBody}
		*msg = messages.SAMessage{_miop}
	case "List":
	default:
		fmt.Println("NamingInvoker:: Operation " + op + " is not implemented by Naming Server")
		os.Exit(0)
	}
}
