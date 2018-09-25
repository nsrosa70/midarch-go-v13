package fibonacciinvoker

import (
	"framework/message"
	"apps/fibonacci/impl"
	"fmt"
)

type FibonacciInvoker struct{}

func (FibonacciInvoker) I_PosInvP(msg *message.Message) {
	op := msg.Payload.(message.MIOP).Body.RequestHeader.Operation
	args := msg.Payload.(message.MIOP).Body.RequestBody.Args

	switch op {
	case "fibo":
		// prepare invocation
		argsX := args.([]interface{})
		p1 := int(argsX[0].(float64))

		//dispatch invocation
		r := impl.Fibo(p1)

		// send reply
		replyHeader := message.ReplyHeader{Status:1} // 1 - Success
		replyBody := message.ReplyBody{Reply: r}

		miopHeader := message.MIOPHeader{Magic:"MIOP"}
		miopBody := message.MIOPBody{ReplyHeader:replyHeader,ReplyBody:replyBody}

		miop := message.MIOP{Header:miopHeader,Body:miopBody}
		*msg = message.Message{miop}
	default:
		fmt.Println("FIBONACCIINVOKER:: Operation "+op+" not supported")
	}
}
