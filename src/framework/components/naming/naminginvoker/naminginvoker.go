package naminginvoker

import (
	"framework/message"
	"fmt"
)

type NamingInvoker struct {}

func (NamingInvoker) I_PosInvP(msg *message.Message){
	fmt.Println(*msg)

	inv := message.Invocation{Op:msg.Payload.(message.MIOP).Body.RequestHeader.Operation,Args:msg.Payload.(message.MIOP).Body.RequestBody.Args}
	msg.Payload = inv
}

func (NamingInvoker) I_PosTerR(msg *message.Message){
	replyHeader := message.ReplyHeader{Status:1}  // 1 - Success
	replyBody := message.ReplyBody{Reply:msg.Payload}

	miopHeader := message.MIOPHeader{Magic:"MIOP"}
	miopBody := message.MIOPBody{ReplyHeader:replyHeader,ReplyBody:replyBody}
	miop := message.MIOP{Header:miopHeader,Body:miopBody}
	msg.Payload = miop
}

