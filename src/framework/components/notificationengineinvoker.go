package components

import (
	"framework/messages"
)

type NotificationEngineInvoker struct{}

var QM NotificationEngine

func (NotificationEngineInvoker) I_PosInvP(msg *messages.SAMessage,r *bool) {
	_op := msg.Payload.(messages.MIOP).Body.RequestHeader.Operation
	_args := msg.Payload.(messages.MIOP).Body.RequestBody.Args
	_argsX := _args.([]interface{})
	*msg = messages.SAMessage{Payload: messages.Invocation{Op: _op, Args: _argsX}}
}

func (NotificationEngineInvoker) I_PosTerR(msg *messages.SAMessage,r *bool) {
	_ter := msg.Payload
	_replyHeader := messages.ReplyHeader{Status: 1} // 1 - Success
	_replyBody := messages.ReplyBody{Reply: _ter}
	_miopHeader := messages.MIOPHeader{Magic: "MIOP"}
	_miopBody := messages.MIOPBody{ReplyHeader: _replyHeader, ReplyBody: _replyBody}
	_miop := messages.MIOP{Header: _miopHeader, Body: _miopBody}
	*msg = messages.SAMessage{Payload: _miop}
}
