package naminginvoker

import (
	"framework/message"
	"framework/components/naming/ior"
	"framework/components/naming/naming"
	"fmt"
	"os"
)

type NamingInvoker struct{}

func (n NamingInvoker) Loop(channels map[string] chan message.Message) {
	var msgPosInvP message.Message
	for {
		select {
		case <-channels["InvP"]:
		case msgPosInvP = <-channels["I_PosInvP"]:
			n.I_PosInvP(&msgPosInvP)
		case channels["TerP"] <- msgPosInvP:
		}
	}
}

func (NamingInvoker) I_PosInvP(msg *message.Message) {
	op := msg.Payload.(message.MIOP).Body.RequestHeader.Operation
	switch op {
	case "Register":
		// Process request
		_args := msg.Payload.(message.MIOP).Body.RequestBody.Args
		_argsX := _args.([]interface{})
		_p1 := _argsX[0].(string)
		_p2 := _argsX[1].(map[string]interface{})
		_ior := ior.IOR{Host: _p2["Host"].(string), Port: int(_p2["Port"].(float64)), Id: int(_p2["Id"].(float64)), Proxy: _p2["Proxy"].(string)}
		_r := naming.Naming{}.Register(_p1, _ior)

		// Send Reply
		_replyHeader := message.ReplyHeader{Status: 1} // 1 - Success
		_replyBody := message.ReplyBody{Reply: _r}
		_miopHeader := message.MIOPHeader{Magic: "MIOP"}
		_miopBody := message.MIOPBody{ReplyHeader: _replyHeader, ReplyBody: _replyBody}
		_miop := message.MIOP{Header: _miopHeader, Body: _miopBody}
		*msg = message.Message{_miop}
	case "Lookup":
		// Process request
		_args := msg.Payload.(message.MIOP).Body.RequestBody.Args
		_argsX := _args.([]interface{})
		_p1 := _argsX[0].(string)
		_r := naming.Naming{}.Lookup(_p1)

		// Send Reply
		_replyHeader := message.ReplyHeader{Status: 1} // 1 - Success
		_replyBody := message.ReplyBody{Reply: _r}
		_miopHeader := message.MIOPHeader{Magic: "MIOP"}
		_miopBody := message.MIOPBody{ReplyHeader: _replyHeader, ReplyBody: _replyBody}
		_miop := message.MIOP{Header: _miopHeader, Body: _miopBody}
		*msg = message.Message{_miop}
	case "List":
	default:
		fmt.Println("NamingInvoker:: Operation " + op + " is not implemented by Naming Server")
		os.Exit(0)
	}
}
