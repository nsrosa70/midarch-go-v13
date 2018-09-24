package naminginvoker

import (
	"framework/message"
	"fmt"
)

type NamingInvoker struct {}

func (NamingInvoker) I_PosInvP(msg *message.Message){
	fmt.Println(*msg)

	inv := message.Invocation{Op:msg.Payload.(message.MIOP).RequestBody.Op,Args:msg.Payload.(message.MIOP).RequestBody.Args}
	msg.Payload = inv
}

func (NamingInvoker) I_PosTerR(msg *message.Message){
	header := message.ReplyHeader{1}  // 1 - Success
	body := message.ReplyBody{Reply:msg.Payload}
	miop := message.MIOP{ReplyHeader:header,ReplyBody:body}
	msg.Payload = miop
}

