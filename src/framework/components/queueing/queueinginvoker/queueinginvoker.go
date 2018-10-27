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
