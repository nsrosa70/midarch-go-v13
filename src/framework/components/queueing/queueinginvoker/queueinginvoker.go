package queueinginvoker

import (
	"framework/message"
	"framework/components/queueing/queueing"
)

type QueueingInvoker struct{}

var QM queueing.QueueingServer

func (QueueingInvoker) I_PosInvP(msg *message.Message) {
	_op := msg.Payload.(message.MIOP).Body.RequestHeader.Operation
	_args := msg.Payload.(message.MIOP).Body.RequestBody.Args
	_argsX := _args.([]interface{})
//	_p1 := _argsX[0].(string)
//	_p2Temp := _argsX[1].(map[string]interface{})
//	_p2 := queueing.MessageMOM{Header: _p2Temp["Header"].(string), PayLoad: _p2Temp["PayLoad"].(string)}
//	_argsInv := []interface{}{_p1, _p2}
	*msg = message.Message{Payload: message.Invocation{Op: _op, Args: _argsX}}
}

func (QueueingInvoker) I_PosTerR(msg *message.Message) {
	_ter := msg.Payload
	_replyHeader := message.ReplyHeader{Status: 1} // 1 - Success
	_replyBody := message.ReplyBody{Reply: _ter}
	_miopHeader := message.MIOPHeader{Magic: "MIOP"}
	_miopBody := message.MIOPBody{ReplyHeader: _replyHeader, ReplyBody: _replyBody}
	_miop := message.MIOP{Header: _miopHeader, Body: _miopBody}
	*msg = message.Message{Payload: _miop}
}
