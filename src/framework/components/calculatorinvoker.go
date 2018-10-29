package components

import (
	"framework/messages"
	"apps/calculator/impl"
)

type CalculatorInvoker struct{}

func (CalculatorInvoker) I_PosInvP(msg *messages.SAMessage) {
	op := msg.Payload.(messages.MIOP).Body.RequestHeader.Operation
	args := msg.Payload.(messages.MIOP).Body.RequestBody.Args

	switch op {
	case "add":
		// prepare invocation
		argsX := args.([]interface{})
		p1 := int(argsX[0].(float64))
		p2 := int(argsX[1].(float64))

		//dispatch invocation
		//fmt.Println("NO PLUGIN")
		r := impl.Add(p1, p2)

		// send reply
		replyHeader := messages.ReplyHeader{Status:1} // 1 - Success
		replyBody := messages.ReplyBody{Reply: r}
		miop := messages.MIOP{Header: messages.MIOPHeader{Magic:"MIOP"},Body: messages.MIOPBody{ReplyHeader:replyHeader,ReplyBody:replyBody}}
		*msg = messages.SAMessage{miop}
	}
}
