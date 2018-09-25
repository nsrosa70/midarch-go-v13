package calculatorinvoker

import (
	"framework/message"
	"apps/calculator/impl"
)

type CalculatorInvoker struct{}

func (CalculatorInvoker) I_PosInvP(msg *message.Message) {
	op := msg.Payload.(message.MIOP).Body.RequestHeader.Operation
	args := msg.Payload.(message.MIOP).Body.RequestBody.Args

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
		replyHeader := message.ReplyHeader{Status:1} // 1 - Success
		replyBody := message.ReplyBody{Reply: r}
		miop := message.MIOP{Header: message.MIOPHeader{Magic:"MIOP"},Body:message.MIOPBody{ReplyHeader:replyHeader,ReplyBody:replyBody}}
		*msg = message.Message{miop}
	}
}
