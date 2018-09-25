package queueinginvoker

import (
	"framework/message"
)

type QueueingInvoker struct {}

func (QueueingInvoker) I_PosInvP(msg *message.Message){
	inv := message.Invocation{Op:msg.Payload.(message.MIOP).Body.RequestHeader.Operation,Args:msg.Payload.(message.MIOP).Body.RequestBody.Args}
	msg.Payload = inv
}

func (QueueingInvoker) I_PosTerR(msg *message.Message){
	replyHeader := message.ReplyHeader{Status:1}  // 1 - Success
	replyBody := message.ReplyBody{Reply:msg.Payload}

	miopHeader := message.MIOPHeader{Magic:"MIOP"}
	miopBody := message.MIOPBody{ReplyHeader:replyHeader,ReplyBody:replyBody}
	miop := message.MIOP{Header:miopHeader,Body:miopBody}
	msg.Payload = miop
}

