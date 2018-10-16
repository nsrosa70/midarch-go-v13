package queueinginvoker

import (
	"framework/message"
	"framework/components/queueing/queueing"
	"fmt"
	"os"
)

type QueueingInvoker struct{}

var QM queueing.Queueing

func (QueueingInvoker) I_PosInvP(msg *message.Message) {
	op := msg.Payload.(message.MIOP).Body.RequestHeader.Operation
	switch op {
	case "Publish":
		// Process request
		_args := msg.Payload.(message.MIOP).Body.RequestBody.Args
		_argsX := _args.([]interface{})
		_p1 := _argsX[0].(string)
		_p2Temp := _argsX[1].(map[string]interface{})
		_p2 := queueing.MessageMOM{Header:_p2Temp["Header"].(string),PayLoad:_p2Temp["PayLoad"].(string)}
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
	case "Subscribe":
	default:
		fmt.Println("Queueing:: Operation " + op + " is not implemented by Queue Server")
		os.Exit(0)
	}
}