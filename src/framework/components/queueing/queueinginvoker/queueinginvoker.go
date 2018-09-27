package queueinginvoker

import (
	"framework/message"
	"framework/components/queueing/queueing"
	"fmt"
	"os"
)

type QueueingInvoker struct{}

func (QueueingInvoker) I_PosInvP(msg *message.Message) {
	op := msg.Payload.(message.MIOP).Body.RequestHeader.Operation
	switch op {
	case "Publish":
		// Process request
		_args := msg.Payload.(message.MIOP).Body.RequestBody.Args
		_argsX := _args.([]interface{})
		_p1 := _argsX[0].(string)
		_p2 := _argsX[1].(string)
		_r := queueing.Queueing{}.Publish(_p1, _p2) // dispatch

		// Send Reply
		_replyHeader := message.ReplyHeader{Status: 1} // 1 - Success
		_replyBody := message.ReplyBody{Reply: _r}
		_miopHeader := message.MIOPHeader{Magic: "MIOP"}
		_miopBody := message.MIOPBody{ReplyHeader: _replyHeader, ReplyBody: _replyBody}
		_miop := message.MIOP{Header: _miopHeader, Body: _miopBody}
		*msg = message.Message{_miop}
	case "Consume":
		_args := msg.Payload.(message.MIOP).Body.RequestBody.Args
		_argsX := _args.([]interface{})
		_p1 := _argsX[0].(string)
		_r := queueing.Queueing{}.Consume(_p1)

		// Send Reply
		_replyHeader := message.ReplyHeader{Status: 1} // 1 - Success
		_replyBody := message.ReplyBody{Reply: _r}
		_miopHeader := message.MIOPHeader{Magic: "MIOP"}
		_miopBody := message.MIOPBody{ReplyHeader: _replyHeader, ReplyBody: _replyBody}
		_miop := message.MIOP{Header: _miopHeader, Body: _miopBody}
		*msg = message.Message{_miop}
	case "subscribe":
	default:
		fmt.Println("Queueing:: Operation " + op + " is not implemented by Queue Server")
		os.Exit(0)
	}
}